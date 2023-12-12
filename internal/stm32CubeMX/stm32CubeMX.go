/*
 * Copyright (c) 2023 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package stm32cubemx

import (
	"bufio"
	"errors"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/open-cmsis-pack/generator-bridge/internal/cbuild"
	"github.com/open-cmsis-pack/generator-bridge/internal/common"
	"github.com/open-cmsis-pack/generator-bridge/internal/utils"
	log "github.com/sirupsen/logrus"
)

func Process(cbuildYmlPath, outPath, cubeMxPath, mxprojectPath string, runCubeMx bool) error {
	var projectFile string
	var parms cbuild.ParamsType

	err := ReadCbuildYmlFile(cbuildYmlPath, outPath, &parms)
	if err != nil {
		return err
	}

	workDir := path.Dir(cbuildYmlPath)
	if parms.OutPath != "" {
		if filepath.IsAbs(parms.OutPath) {
			workDir = parms.OutPath
		} else {
			workDir = path.Join(workDir, parms.OutPath)
		}
	} else {
		workDir = path.Join(workDir, outPath)
	}
	workDir = filepath.Clean(workDir)
	workDir = filepath.ToSlash(workDir)

	err = os.MkdirAll(workDir, os.ModePerm)
	if err != nil {
		return err
	}

	if runCubeMx {
		cubeIocPath := workDir
		lastPath := filepath.Base(cubeIocPath)
		if lastPath != "STM32CubeMX" {
			cubeIocPath = path.Join(cubeIocPath, "STM32CubeMX")
		}
		cubeIocPath = path.Join(cubeIocPath, "STM32CubeMX.ioc")

		if utils.FileExists(cubeIocPath) {
			err := Launch(cubeIocPath, "")
			if err != nil {
				return err
			}
		} else {
			projectFile, err = WriteProjectFile(workDir, &parms)
			if err != nil {
				return nil
			}
			log.Infof("Generated file: %v", projectFile)

			err := Launch("", projectFile)
			if err != nil {
				return err
			}
		}

		err = ReadContexts(cubeIocPath)
		if err != nil {
			return err
		}

		tmpPath, _ := filepath.Split(cubeIocPath)
		mxprojectPath = path.Join(tmpPath, ".mxproject")
	}
	mxproject, err := IniReader(mxprojectPath, false)
	if err != nil {
		return err
	}

	err = WriteCgenYml(workDir, mxproject, parms)
	if err != nil {
		return err
	}

	return nil
}

func Launch(iocFile, projectFile string) error {
	log.Infof("Launching STM32CubeMX...")

	const cubeEnvVar = "STM32CubeMX_PATH"
	cubeEnv := os.Getenv(cubeEnvVar)
	if cubeEnv == "" {
		return errors.New("environment variable for CubeMX not set: " + cubeEnvVar)
	}

	pathJava := path.Join(cubeEnv, "jre", "bin", "java.exe")
	pathCubeMx := path.Join(cubeEnv, "STM32CubeMX.exe")

	var cmd *exec.Cmd
	if iocFile != "" {
		cmd = exec.Command(pathJava, "-jar", pathCubeMx, iocFile)
	} else if projectFile != "" {
		cmd = exec.Command(pathJava, "-jar", pathCubeMx, "-s", projectFile)
	} else {
		cmd = exec.Command(pathJava, "-jar", pathCubeMx)
	}
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func WriteProjectFile(workDir string, parms *cbuild.ParamsType) (string, error) {
	filePath := filepath.Join(workDir, "project.script")
	log.Infof("Writing CubeMX project file %v", filePath)

	var text utils.TextBuilder
	if parms.Board != "" {
		text.AddLine("loadboard", parms.Board, "allmodes")
	} else {
		text.AddLine("load", parms.Device)
	}
	text.AddLine("project name", "STM32CubeMX")
	text.AddLine("project toolchain", utils.AddQuotes("MDK-ARM V5"))

	cubeWorkDir := workDir
	if runtime.GOOS == "windows" {
		cubeWorkDir = filepath.FromSlash(cubeWorkDir)
	}
	text.AddLine("project path", utils.AddQuotes(cubeWorkDir))
	text.AddLine("SetCopyLibrary", utils.AddQuotes("copy only"))

	if utils.FileExists(filePath) {
		os.Remove(filePath)
	}

	err := os.WriteFile(filePath, []byte(text.GetLine()), 0600)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func ReadCbuildYmlFile(path, outPath string, parms *cbuild.ParamsType) error {
	log.Infof("Reading cbuild.yml file: '%v'", path)
	err := cbuild.Read(path, outPath, parms)
	if err != nil {
		return err
	}

	return nil
}

var filterFiles = map[string]string{
	"system_":   "system_ file (delivered from elsewhere)",
	"Templates": "Templates file (mostly not present)",
}

func FilterFile(file string) bool {
	for key, value := range filterFiles {
		if strings.Contains(file, key) {
			log.Infof("ignoring %v: %v", value, file)
			return true
		}
	}

	return false
}

func FindMxProject(subsystem *cbuild.SubsystemType, mxprojectAll MxprojectAllType) (MxprojectType, error) {
	if len(mxprojectAll.Mxproject) == 0 {
		return MxprojectType{}, errors.New("no .mxproject read")
	} else if len(mxprojectAll.Mxproject) == 1 {
		mxproject := mxprojectAll.Mxproject[0]
		return mxproject, nil
	}

	coreName := subsystem.CoreName
	trustzone := subsystem.TrustZone
	if trustzone == "off" {
		trustzone = ""
	}
	for id := range mxprojectAll.Mxproject {
		mxproject := mxprojectAll.Mxproject[id]
		if mxproject.CoreName == coreName && mxproject.Trustzone == trustzone {
			return mxproject, nil
		}
	}

	return MxprojectType{}, nil
}

type PinDefinition struct {
	p			string
	pin			string
	port		string
	mode		string
	pull 		string
	speed		string
	alternate	string
}

func ReadContexts(iocFile string) error {
	contextMap, err := createContextMap(iocFile)
	if err != nil {
		return err
	}

	contexts, err := getContexts(contextMap)
	if err != nil {
		return err
	}

	deviceFamily, err := getDeviceFamily(contextMap)
	if err != nil {
		return err
	}

	workDir := path.Dir(iocFile)
	workDirAbs, err := filepath.Abs(workDir)
	if err != nil {
		return err
	}

	mspName := deviceFamily + "xx_hal_msp.c"
//	pinConfigMap, err := createPinConfigMap(mspName)
	mspFolder := contextMap["ProjectManager"]["MainLocation"]
	if mspFolder == "" {
		return errors.New("main location missing")
	}

	projectIndex := 1
	if len(contexts) == 0 {
		msp := path.Join(workDirAbs, mspFolder, mspName)
		fMsp, err := os.Open(msp)
		if err != nil {
			return err
		}
		defer fMsp.Close()

		fName := "MX_Device.h"
		fPath := path.Join(path.Dir(workDir), "drv_cfg", "cproject_" + strconv.Itoa(projectIndex))
		if _, err := os.Stat(fPath); err != nil {
			os.MkdirAll(fPath, 0750)
		}
		fPath = path.Join(fPath, fName)
		fMxDevice, err := os.Create(fPath)
		if err != nil {
			return err
		}
		defer fMxDevice.Close()

		mxDeviceWriteHeader(fMxDevice, fName)
        peripherals, err := getPeripherals1(contextMap)
		if err != nil {
			return err
		}
        for _, peripheral := range(peripherals) {
            vmode := getVirtualMode(contextMap, peripheral)
            pins := getPins(contextMap, fMsp, peripheral)
            err := mxDeviceWritePeripheralCfg(fMxDevice, peripheral, vmode, pins)
			if err != nil {
				return err
			}
		}
        fMxDevice.WriteString("\n#endif  /* __MX_DEVICE_H */\n")
	} else {
		CONTEXT := make(map[string]string)
		CONTEXT["CortexM33S"] = "Secure"
		CONTEXT["CortexM33NS"] = "NonSecure"
		CONTEXT["CortexM4"] = "CM4"
		CONTEXT["CortexM7"] = "CM7"
		for _, context := range contexts {
			contextFolder := CONTEXT[context]
			if contextFolder == "" {
				print("Cannot find ", mspName)
				return errors.New("Cannot find " + mspName)
			}
			msp := path.Join(workDirAbs, contextFolder, mspFolder, mspName)
			fMsp, err := os.Open(msp)
			if err != nil {
				return err
			}
	
			fName := "MX_Device.h"
			fPath := path.Join(path.Dir(workDir), "drv_cfg", "cproject_" + strconv.Itoa(projectIndex))
			if _, err := os.Stat(fPath); err != nil {
				os.MkdirAll(fPath, 0750)
			}
			fPath = path.Join(fPath, fName)
			fMxDevice, err := os.Create(fPath)
			if err != nil {
				return err
			}

			projectIndex += 1
	
			mxDeviceWriteHeader(fMxDevice, fName)
			peripherals, err := getPeripherals(contextMap, context)
			if err != nil {
				fMxDevice.Close()
				return err
			}
			for _, peripheral := range(peripherals) {
				vmode := getVirtualMode(contextMap, peripheral)
				pins := getPins(contextMap, fMsp, peripheral)
				err := mxDeviceWritePeripheralCfg(fMxDevice, peripheral, vmode, pins)
				if err != nil {
					fMxDevice.Close()
					return err
				}
			}
			fMxDevice.WriteString("\n#endif  /* __MX_DEVICE_H */\n")
			fMxDevice.Close()
			}
	}
	return nil
}

func createContextMap(iocFile string) (map[string]map[string]string, error) {
	contextMap := make(map[string]map[string]string)
	var xMap map[string]string

	fIoc, err := os.Open(iocFile)
	if err != nil {
		return nil, err
	}
	defer fIoc.Close()
	iocScan := bufio.NewScanner(fIoc)
	iocScan.Split(bufio.ScanLines)
	for iocScan.Scan() {
		split := strings.Split(iocScan.Text(), "=")
		if len(split) > 1 {
			leftParts := strings.Split(split[0], ".")
			if len(leftParts) > 1 {
				if contextMap[leftParts[0]] == nil {
					xMap = make(map[string]string)
					contextMap[leftParts[0]] = xMap
				}
				xMap[leftParts[1]] = split[1]
			}
		}
	}
	return contextMap, nil
}

/*
func createPinConfigMap(mspName string) (map[string]PinDefinition, error) {
	pinConfigMap := make(map[string]PinDefinition)
	var instance string
	fIoc, err := os.Open(mspName)
	if err != nil {
		return nil, err
	}
	defer fIoc.Close()
	mspScan := bufio.NewScanner(fIoc)
	mspScan.Split(bufio.ScanLines)
	s := "->Instance=="
	h := "HAL_GPIO_Init"
	for mspScan.Scan() {
		line := mspScan.Text()
		if line == "}" {				// end of function
			instance = ""				// reset instance
		}
		if len(instance) == 0 {			// no instance yet
			idx := strings.Index(line, s)
			if idx != -1 {
				inst := strings.Split(line[idx:], ")")[0]
				if len(inst) > 0 {
					instance = inst
				}
			}
		} else {						// there was an instance
			idx := strings.Index(line, h)
			if idx != -1 {
				pinConfigMap[instance] = pinDef
			}
		}
	}
	return pinConfigMap, nil
}
 */
func getContexts(contextMap map[string]map[string]string) ([]string, error) {
	var contexts [] string
	head := contextMap["Mcu"]
	if len(head) > 0 {
		for key, content := range head {
			if strings.HasPrefix(key, "Context") {
	 			l := len(key)
	 			if l > 0 && key[l-1] >= '0' && key[l-1] <= '9' {
					contexts = append(contexts, content)
				}
			}
		}
	}
	return contexts, nil
}

func getDeviceFamily(contextMap map[string]map[string]string) (string, error) {
	family := contextMap["Mcu"]["Family"]
	if family != "" {
		if strings.HasPrefix(family, "STM32") {
			return family, nil
		}
	}
	return "", errors.New("missing device family")
}

func getPeripherals1(contextMap map[string]map[string]string) ([]string, error) {
	PERIPHERALS := [...]string{"USART", "UART", "LPUART", "SPI", "I2C", "ETH", "SDMMC", "CAN", "USB", "SDIO", "FDCAN"}
	var peripherals [] string
	mcu := contextMap["Mcu"]
	if mcu != nil {
		for ip, peri := range(mcu) {
			if strings.HasPrefix(ip, "IP") {
				for _, peripheral := range(PERIPHERALS) {
					if strings.HasPrefix(peri, peripheral) {
						peripherals = append(peripherals, peri)
					}
				}
			}
		}
	} else {
		return nil, errors.New("peripheral not found in Mcu")
	}
	return peripherals, nil
}

func getPeripherals(contextMap map[string]map[string]string, context string) ([]string, error) {
	PERIPHERALS := [...]string{"USART", "UART", "LPUART", "SPI", "I2C", "ETH", "SDMMC", "CAN", "USB", "SDIO", "FDCAN"}
	var peripherals [] string
	contextIps := contextMap[context]
	if contextIps == nil {
		return nil, errors.New("context not found in ioc")
	}
	contextIpsLine := contextIps["IPs"]
	if len(contextIpsLine) == 0 {
		return nil, errors.New("IPs not found in context")
	}
	mcu := contextMap["Mcu"]
	if mcu != nil {
		for ip, peri := range(mcu) {
			if strings.HasPrefix(ip, "IP") {
				if strings.Contains(contextIpsLine, peri) {
					for _, peripheral := range(PERIPHERALS) {
						if strings.HasPrefix(peri, peripheral) {
							peripherals = append(peripherals, peri)
						}
					}
				}
			}
		}
	} else {
		return nil, errors.New("peripheral not found in Mcu")
	}
	return peripherals, nil
}

func getVirtualMode(contextMap map[string]map[string]string, peripheral string) string {
	peri := contextMap[peripheral]
	if len(peri) > 0 {
		for vm, vmode := range(peri) {
			if strings.HasPrefix(vm, "VirtualMode") {
				return vmode
			}
		}
	}
	return ""
}

func getPins(contextMap map[string]map[string]string, fMsp *os.File, peripheral string) map[string]PinDefinition {
	pinsName := make(map[string]string)
	pinsLabel := make(map[string]string)
	pinsInfo := make(map[string]PinDefinition)
	for key, signal := range(contextMap) {
		if !strings.HasPrefix(key, "VP") {
			peri := signal["Signal"]
			if strings.HasPrefix(peri, peripheral) {
				pinsName[key] = peri
				label := signal["GPIO_Label"]
				if len(label) > 0 {
					label = strings.Split(label, " ")[0]
					label = replaceSpecialChars(label, "_")
					pinsLabel[key] = strings.ReplaceAll(label, ".", "_")
				}
			}
		}
	}
	for pin, name := range(pinsName) {
		p := strings.Split(pin, "\\")[0]
		p = strings.Split(p, "(")[0]
		p = strings.Split(p, " ")[0]
		p = strings.Split(p, "_")[0]
		p = strings.Split(p, "-")[0]
		label := pinsLabel[pin]
		info := getPinConfiguration(fMsp, peripheral, p, label)
		pinsInfo[name] = info
	}
	return pinsInfo
}

func replaceSpecialChars(label string, ch string) string {
	specialCharacter := []string{"!", "@", "#", "$", "%", "^", "&", "*", "(",  "+", "=", "-", "_", "[", "]", "{", "}",
								 ";", ":", ",", ".", "?", "/", "\\", "|", "~", "`", "\"", "'", "<", ">", " "}
	for _, spec := range(specialCharacter) {
		label = strings.ReplaceAll(label, spec, ch)
	}
	return label
}

func getDigitAtEnd(pin string) string {
	re := regexp.MustCompile("[0-9]+$")
	numbers := re.FindAllString(pin, -1)
	if numbers != nil {
		return numbers[0]
	}
	return ""
}

func getPinConfiguration(fMsp *os.File, peripheral string, pin string, label string) PinDefinition {
	var pinInfo	PinDefinition

	pinNum := getDigitAtEnd(pin)
	gpioPin := "GPIO_PIN_" + pinNum
	port := strings.Split(strings.Split(pin, "P")[1], pinNum)[0]
	gpioPort := "GPIO" + port

	section := false
	fMsp.Seek(0, 0)
	mspScan := bufio.NewScanner(fMsp)
	mspScan.Split(bufio.ScanLines)
	s := "->Instance=="
	h := "HAL_GPIO_Init"
	addLine := false
	value := ""
	for mspScan.Scan() {
		line := mspScan.Text()
		if line == "}" {				// end of function
			section = false				// reset instance
		}
		if strings.Contains(line, s) && strings.Contains(line, peripheral) {
			section = true
		}
		if section {
			if strings.Contains(line, h) {
				if strings.Contains(line, gpioPort) || strings.Contains(line, label + "_GPIO_Port") {
					values := strings.Split(pinInfo.pin, "|")
					for _, val := range(values) {
						val = strings.TrimRight(strings.TrimLeft(val, " "), " ")
						if val == gpioPin || val == (label + "_Pin") {
							pinInfo.p = pin;
							pinInfo.pin = gpioPin
							pinInfo.port = gpioPort
							return pinInfo
						}
					}
				}
			}
			if addLine {
				value += strings.TrimLeft(line, " ")
				if strings.Contains(value, ";") {
					pinInfo.pin = strings.Split(value, ";")[0]
					addLine = false
				}
			} else {
				assign := strings.Split(line, "=")
				if len(assign) > 1 {
					left := assign[0]
					value = strings.TrimLeft(assign[1], " ")
					switch {
					case strings.Contains(left, ".Pin"):
						if strings.Contains(value, ";") {
							pinInfo.pin = strings.Split(value, ";")[0]
						} else {
							addLine = true
						}
					case strings.Contains(left, ".Port"):
						pinInfo.port = strings.Split(value, ";")[0]
					case strings.Contains(left, ".Mode"):
						pinInfo.mode = strings.Split(value, ";")[0]
					case strings.Contains(left, ".Pull"):
						pinInfo.pull = strings.Split(value, ";")[0]
					case strings.Contains(left, ".Speed"):
						pinInfo.speed = strings.Split(value, ";")[0]
					case strings.Contains(left, ".Alternate"):
						pinInfo.alternate = strings.Split(value, ";")[0]
					}
				}
			}
		}
	}
	return PinDefinition{}
}

func mxDeviceWriteHeader(fMxDevice *os.File, fName string) error {
	now := time.Now()
	dtString := now.Format("02/01/2006 15:04:05")

	fMxDevice.WriteString("/******************************************************************************\n")
	fMxDevice.WriteString(" * File Name   : " + fName + "\n")
	fMxDevice.WriteString(" * Date        : " + dtString + "\n")
	fMxDevice.WriteString(" * Description : STM32Cube MX parameter definitions\n")
	fMxDevice.WriteString(" * Note        : This file is generated with a generator out of the\n")
	fMxDevice.WriteString(" *               STM32CubeMX project and its generated files (DO NOT EDIT!)\n")
	fMxDevice.WriteString(" ******************************************************************************/\n\n")
	fMxDevice.WriteString("#ifndef __MX_DEVICE_H\n")
	fMxDevice.WriteString("#define __MX_DEVICE_H\n\n")
	return nil
}

func mxDeviceWritePeripheralCfg(fMxDevice *os.File, peripheral string, vmode string, pins map[string]PinDefinition) error {
	str := "\n/*------------------------------ " + peripheral
	if len(str) < 49 {
		str += strings.Repeat(" ", 49-len(str))
	}
	str += "-----------------------------*/\n"
	fMxDevice.WriteString(str)
	fMxDevice.WriteString(createDefine(peripheral, "1") + "\n\n")
	if vmode != "" {
		fMxDevice.WriteString("/* Virtual mode */\n")
		fMxDevice.WriteString(createDefine(peripheral + "_VM", vmode) + "\n")
		fMxDevice.WriteString(createDefine(peripheral + "_" + vmode, "1") + "\n\n")
	}
	if len(pins) != 0 {
		fMxDevice.WriteString("/* Pins */\n")
		for pin, pinDef := range(pins) {
			fMxDevice.WriteString("\n/* " + pin + " */\n")
			if len(pinDef.p) != 0 {
				fMxDevice.WriteString(createDefine(pin + "_Pin", pinDef.p) + "\n")
			}
			if len(pinDef.pin) != 0 {
				fMxDevice.WriteString(createDefine(pin + "_GPIO_Pin", pinDef.pin) + "\n")
			}
			if len(pinDef.port) != 0 {
				fMxDevice.WriteString(createDefine(pin + "_GPIOx", pinDef.port) + "\n")
			}
			if len(pinDef.mode) != 0 {
				fMxDevice.WriteString(createDefine(pin + "_GPIO_Mode", pinDef.mode) + "\n")
			}
			if len(pinDef.pull) != 0 {
				fMxDevice.WriteString(createDefine(pin + "_GPIO_PuPd", pinDef.pull) + "\n")
			}
			if len(pinDef.speed) != 0 {
				fMxDevice.WriteString(createDefine(pin + "_GPIO_Speed", pinDef.speed) + "\n")
			}
			if len(pinDef.alternate) != 0 {
				fMxDevice.WriteString(createDefine(pin + "_GPIO_AF", pinDef.alternate) + "\n")
			}
		}
	}

	return nil
}

func createDefine(name string, value string) string {
	invalidChars := [...]string{"=", " ", "/", "(", ")", "[", "]", "\\", "-"}

	for _, ch := range(invalidChars) {
		name = strings.ReplaceAll(name, ch, "_")
		value = strings.ReplaceAll(value, ch, "_")
	}
	name = "MX_" + name
	if len(name) < 39 {
		name += strings.Repeat(" ", 39-len(name))
	}
	define := "#define " + name + value
	return define
}

func WriteCgenYml(outPath string, mxprojectAll MxprojectAllType, inParms cbuild.ParamsType) error {
	for id := range inParms.Subsystem {
		subsystem := &inParms.Subsystem[id]

		mxproject, err := FindMxProject(subsystem, mxprojectAll)
		if err != nil {
			continue
		}
		err = WriteCgenYmlSub(outPath, mxproject, subsystem)
		if err != nil {
			return err
		}
	}

	return nil
}

func WriteCgenYmlSub(outPath string, mxproject MxprojectType, subsystem *cbuild.SubsystemType) error {
	outName := subsystem.SubsystemIdx.Project + ".cgen.yml"
	outFile := path.Join(outPath, outName)
	var cgen cbuild.CgenType

	lastPath := filepath.Base(outPath)
	var relativePathAdd string
	if lastPath != "STM32CubeMX" {
		relativePathAdd = path.Join(relativePathAdd, "STM32CubeMX")
	}
	relativePathAdd = path.Join(relativePathAdd, "MDK-ARM")

	cgen.GeneratorImport.ForBoard = subsystem.Board
	cgen.GeneratorImport.ForDevice = subsystem.Device
	cgen.GeneratorImport.Define = append(cgen.GeneratorImport.Define, mxproject.PreviousUsedKeilFiles.CDefines...)

	for id := range mxproject.PreviousUsedKeilFiles.HeaderPath {
		headerPath := mxproject.PreviousUsedKeilFiles.HeaderPath[id]
		headerPath, _ = utils.ConvertFilename(outPath, headerPath, relativePathAdd)
		cgen.GeneratorImport.AddPath = append(cgen.GeneratorImport.AddPath, headerPath)
	}

	var groupSrc cbuild.CgenGroupsType
	var groupHalDriver cbuild.CgenGroupsType
	var groupTz cbuild.CgenGroupsType

	groupSrc.Group = "CubeMX"
	groupHalDriver.Group = "HAL Driver"
	groupHalFilter := "HAL_Driver"

	for id := range mxproject.PreviousUsedKeilFiles.SourceFiles {
		file := mxproject.PreviousUsedKeilFiles.SourceFiles[id]
		if FilterFile(file) {
			continue
		}
		file, _ = utils.ConvertFilename(outPath, file, relativePathAdd)

		if strings.Contains(file, groupHalFilter) {
			var cgenFile cbuild.CgenFilesType
			cgenFile.File = file
			groupHalDriver.Files = append(groupHalDriver.Files, cgenFile)
		} else {
			var cgenFile cbuild.CgenFilesType
			cgenFile.File = file
			groupSrc.Files = append(groupSrc.Files, cgenFile)
		}
	}

	cgen.GeneratorImport.Groups = append(cgen.GeneratorImport.Groups, groupSrc)
	cgen.GeneratorImport.Groups = append(cgen.GeneratorImport.Groups, groupHalDriver)

	if subsystem.TrustZone == "non-secure" {
		groupTz.Group = "CMSE Library"
		var cgenFile cbuild.CgenFilesType
		cgenFile.File = "$cmse-lib("
		cgenFile.File += subsystem.SubsystemIdx.SecureContextName
		cgenFile.File += ")$"
		groupTz.Files = append(groupTz.Files, cgenFile)
		cgen.GeneratorImport.Groups = append(cgen.GeneratorImport.Groups, groupTz)
	}

	return common.WriteYml(outFile, &cgen)
}
