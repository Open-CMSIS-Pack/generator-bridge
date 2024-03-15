/*
 * Copyright (c) 2023-2024 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package stm32cubemx

import (
	"bufio"
	"errors"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/open-cmsis-pack/generator-bridge/internal/cbuild"
)

type PinDefinition struct {
	p         string
	pin       string
	port      string
	mode      string
	pull      string
	speed     string
	alternate string
}

func ReadContexts(iocFile string, params cbuild.ParamsType) error {
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

	mainFolder := contextMap["ProjectManager"]["MainLocation"]
	if mainFolder == "" {
		return errors.New("main location missing")
	}
	mspName := deviceFamily + "xx_hal_msp.c"

	var cfgPath string
	if len(contexts) == 0 {
		cfgPath = path.Join("drv_cfg", params.Subsystem[0].SubsystemIdx.Project)
		err := writeMXdeviceH(contextMap, workDir, mainFolder, mspName, cfgPath, "", params)
		if err != nil {
			return err
		}
	} else {
		for _, context := range contexts {
			var coreName string
			var contextFolder string
			var projectPart string
			re := regexp.MustCompile("[0-9]+")
			coreNameNumbers := re.FindAllString(context, -1)
			if len(coreNameNumbers) == 1 {
				coreName = "Cortex-M" + coreNameNumbers[0]
				contextFolder = "CM" + coreNameNumbers[0]
				projectPart = contextFolder
			}

			var trustzone string
			contextLen := len(context)
			if contextLen > 0 {
				if strings.LastIndex(context, "S") == contextLen-1 {
					if strings.LastIndex(context, "NS") == contextLen-2 {
						trustzone = "NonSecure"
						projectPart = "non-secure"
					} else {
						trustzone = "Secure"
						projectPart = "secure"
					}
					contextFolder = trustzone
				}
			}

			if len(contextFolder) == 0 {
				return errors.New("Cannot find context " + context)
			}
			for _, subsystem := range params.Subsystem {
				if subsystem.CoreName == coreName {
					if len(subsystem.TrustZone) == 0 {
						cfgPath = path.Join("drv_cfg", subsystem.SubsystemIdx.Project)
						break
					}
					if subsystem.TrustZone == projectPart {
						cfgPath = path.Join("drv_cfg", subsystem.SubsystemIdx.Project)
						break
					}
				}
			}
			err := writeMXdeviceH(contextMap, workDir, path.Join(contextFolder, mainFolder), mspName, cfgPath, context, params)
			if err != nil {
				return err
			}
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

func writeMXdeviceH(contextMap map[string]map[string]string, workDir string, mainFolder string, mspName string, cfgPath string, context string, params cbuild.ParamsType) error {
	_ = params
	workDirAbs, err := filepath.Abs(workDir)
	if err != nil {
		return err
	}

	main := path.Join(workDirAbs, mainFolder, "main.c")
	fMain, err := os.Open(main)
	if err != nil {
		return err
	}
	defer fMain.Close()

	msp := path.Join(workDirAbs, mainFolder, mspName)
	fMsp, err := os.Open(msp)
	if err != nil {
		return err
	}
	defer fMsp.Close()

	fName := "MX_Device.h"
	fPath := path.Join(path.Dir(workDir), cfgPath)
	if _, err := os.Stat(fPath); err != nil {
		err = os.MkdirAll(fPath, 0750)
		if err != nil {
			return err
		}
	}
	fPath = path.Join(fPath, fName)
	fMxDevice, err := os.Create(fPath)
	if err != nil {
		return err
	}
	defer fMxDevice.Close()

	out := bufio.NewWriter(fMxDevice)
	defer out.Flush()
	err = mxDeviceWriteHeader(out, fName)
	if err != nil {
		return err
	}
	peripherals, err := getPeripherals(contextMap, context)
	if err != nil {
		return err
	}
	for _, peripheral := range peripherals {
		vmode := getVirtualMode(contextMap, peripheral)
		i2cInfo, err := getI2cInfo(fMain, peripheral)
		if err != nil {
			return err
		}
		pins, err := getPins(contextMap, fMsp, peripheral)
		if err != nil {
			return err
		}
		err = mxDeviceWritePeripheralCfg(out, peripheral, vmode, i2cInfo, pins)
		if err != nil {
			return err
		}
	}
	_, err = out.WriteString("\n#endif  /* __MX_DEVICE_H */\n")
	if err != nil {
		return err
	}
	return nil
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
func getContexts(contextMap map[string]map[string]string) (map[int]string, error) {
	contexts := make(map[int]string)
	head := contextMap["Mcu"]
	if len(head) > 0 {
		for key, content := range head {
			if strings.HasPrefix(key, "Context") {
				l := len(key)
				if l > 0 && key[l-1] >= '0' && key[l-1] <= '9' {
					i, err := strconv.Atoi(string(key[l-1]))
					if err != nil {
						return nil, err
					}
					contexts[i] = content
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

func getPeripherals(contextMap map[string]map[string]string, context string) ([]string, error) {
	PERIPHERALS := [...]string{"USART", "UART", "LPUART", "SPI", "I2C", "ETH", "SDMMC", "CAN", "USB", "SDIO", "FDCAN"}
	var peripherals []string
	var contextIpsLine string
	if len(context) > 0 {
		contextIps, ok := contextMap[context]
		if !ok {
			return nil, errors.New("context not found in ioc")
		}
		contextIpsLine, ok = contextIps["IPs"]
		if !ok {
			return nil, errors.New("IPs not found in context")
		}
	}
	mcu := contextMap["Mcu"]
	if mcu != nil {
		for ip, peri := range mcu {
			if strings.HasPrefix(ip, "IP") {
				if len(context) == 0 || strings.Contains(contextIpsLine, peri) {
					for _, peripheral := range PERIPHERALS {
						if strings.HasPrefix(peri, peripheral) {
							peripherals = append(peripherals, peri)
							break
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
		for vm, vmode := range peri {
			if strings.HasPrefix(vm, "VirtualMode") {
				return vmode
			}
		}
	}
	return ""
}

func getPins(contextMap map[string]map[string]string, fMsp *os.File, peripheral string) (map[string]PinDefinition, error) {
	pinsName := make(map[string]string)
	pinsLabel := make(map[string]string)
	pinsInfo := make(map[string]PinDefinition)
	for key, signal := range contextMap {
		if !strings.HasPrefix(key, "VP") {
			peri := signal["Signal"]
			if strings.HasPrefix(peri, peripheral) {
				pinsName[key] = peri
				label, ok := signal["GPIO_Label"]
				if ok {
					label = strings.Split(label, " ")[0]
					label = replaceSpecialChars(label, "_")
					pinsLabel[key] = strings.ReplaceAll(label, ".", "_")
				}
			}
		}
	}
	for pin, name := range pinsName {
		p := strings.Split(pin, "\\")[0]
		p = strings.Split(p, "(")[0]
		p = strings.Split(p, " ")[0]
		p = strings.Split(p, "_")[0]
		p = strings.Split(p, "-")[0]
		label := pinsLabel[pin]
		info, err := getPinConfiguration(fMsp, peripheral, p, label)
		if err != nil {
			return nil, err
		}
		pinsInfo[name] = info
	}
	return pinsInfo, nil
}

func replaceSpecialChars(label string, ch string) string {
	specialCharacter := [...]string{"!", "@", "#", "$", "%", "^", "&", "*", "(", "+", "=", "-", "_", "[", "]", "{", "}",
		";", ":", ",", ".", "?", "/", "\\", "|", "~", "`", "\"", "'", "<", ">", " "}
	for _, spec := range specialCharacter {
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

// Get i2c info (filter, coefficients)
func getI2cInfo(fMain *os.File, peripheral string) (map[string]string, error) {
	info := make(map[string]string)
	if strings.HasPrefix(peripheral, "I2C") {
		_, err := fMain.Seek(0, 0)
		if err != nil {
			return nil, err
		}
		section := false

		mainScan := bufio.NewScanner(fMain)
		mainScan.Split(bufio.ScanLines)
		for mainScan.Scan() {
			line := mainScan.Text()
			if !section {
				if strings.HasPrefix(line, "static void MX_"+peripheral+"_Init") && !strings.Contains(line, ";") {
					section = true // Start of section: static void MX_I2Cx_Init
				}
			} else { // Parse section: static void MX_I2Cx_Init
				if strings.HasPrefix(line, "}") {
					break // End of section: static void MX_I2Cx_Init
				}
				if strings.Contains(line, "HAL_I2CEx_ConfigAnalogFilter") {
					if strings.Contains(line, "I2C_ANALOGFILTER_ENABLE") {
						info["ANF_ENABLE"] = "1"
					} else {
						info["ANF_ENABLE"] = "0"
					}
				}
				if strings.Contains(line, "HAL_I2CEx_ConfigDigitalFilter") {
					dnf := strings.Split(strings.Split(line, ",")[1], ")")[0]
					dnf = strings.TrimRight(strings.TrimLeft(dnf, "\t "), "\t ")
					info["DNF"] = dnf
				}
			}
		}
	}
	return info, nil
}

func getPinConfiguration(fMsp *os.File, peripheral string, pin string, label string) (PinDefinition, error) {
	var pinInfo PinDefinition

	pinNum := getDigitAtEnd(pin)
	gpioPin := "GPIO_PIN_" + pinNum
	port := strings.Split(strings.Split(pin, "P")[1], pinNum)[0]
	gpioPort := "GPIO" + port

	section := false
	_, err := fMsp.Seek(0, 0)
	if err != nil {
		return PinDefinition{}, err
	}
	mspScan := bufio.NewScanner(fMsp)
	mspScan.Split(bufio.ScanLines)
	s := "->Instance=="
	h := "HAL_GPIO_Init"
	addLine := false
	value := ""
	for mspScan.Scan() {
		line := mspScan.Text()
		if line == "}" { // end of function
			section = false // reset instance
		}
		if strings.Contains(line, s) && strings.Contains(line, peripheral) {
			section = true
		}
		if section {
			if strings.Contains(line, h) {
				if strings.Contains(line, gpioPort) || strings.Contains(line, label+"_GPIO_Port") {
					values := strings.Split(pinInfo.pin, "|")
					for _, val := range values {
						val = strings.TrimRight(strings.TrimLeft(val, "\t "), "\t ")
						if val == gpioPin || val == (label+"_Pin") {
							pinInfo.p = pin
							pinInfo.pin = gpioPin
							pinInfo.port = gpioPort
							return pinInfo, nil
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
	return PinDefinition{}, nil
}

func mxDeviceWriteHeader(out *bufio.Writer, fName string) error {
	now := time.Now()
	dtString := now.Format("02/01/2006 15:04:05")

	var err error

	if _, err = out.WriteString("/******************************************************************************\n"); err != nil {
		return err
	}
	if _, err = out.WriteString(" * File Name   : " + fName + "\n"); err != nil {
		return err
	}
	if _, err = out.WriteString(" * Date        : " + dtString + "\n"); err != nil {
		return err
	}
	if _, err = out.WriteString(" * Description : STM32Cube MX parameter definitions\n"); err != nil {
		return err
	}
	if _, err = out.WriteString(" * Note        : This file is generated with a generator out of the\n"); err != nil {
		return err
	}
	if _, err = out.WriteString(" *               STM32CubeMX project and its generated files (DO NOT EDIT!)\n"); err != nil {
		return err
	}
	if _, err = out.WriteString(" ******************************************************************************/\n\n"); err != nil {
		return err
	}
	if _, err = out.WriteString("#ifndef __MX_DEVICE_H\n"); err != nil {
		return err
	}
	_, err = out.WriteString("#define __MX_DEVICE_H\n\n")
	return err
}

func mxDeviceWritePeripheralCfg(out *bufio.Writer, peripheral string, vmode string, i2cInfo map[string]string, pins map[string]PinDefinition) error {
	var err error

	str := "\n/*------------------------------ " + peripheral
	if len(str) < 49 {
		str += strings.Repeat(" ", 49-len(str))
	}
	str += "-----------------------------*/\n"
	if _, err = out.WriteString(str); err != nil {
		return err
	}
	if err = writeDefine(out, peripheral, "1\n"); err != nil {
		return err
	}
	if i2cInfo != nil {
		if _, err = out.WriteString("/* Filter Settings */\n"); err != nil {
			return err
		}
		for item, value := range i2cInfo {
			if err = writeDefine(out, peripheral+"_"+item, value); err != nil {
				return err
			}
		}
		if _, err = out.WriteString("\n"); err != nil {
			return err
		}
	}
	if vmode != "" {
		if _, err = out.WriteString("/* Virtual mode */\n"); err != nil {
			return err
		}
		if err = writeDefine(out, peripheral+"_VM", vmode); err != nil {
			return err
		}
		if err = writeDefine(out, peripheral+"_"+vmode, "1"); err != nil {
			return err
		}
	}
	if len(pins) != 0 {
		if _, err = out.WriteString("/* Pins */\n"); err != nil {
			return err
		}
		for pin, pinDef := range pins {
			if _, err = out.WriteString("\n/* " + pin + " */\n"); err != nil {
				return err
			}
			if err = writeDefine(out, pin+"_Pin", pinDef.p); err != nil {
				return err
			}
			if err = writeDefine(out, pin+"_GPIO_Pin", pinDef.pin); err != nil {
				return err
			}
			if err = writeDefine(out, pin+"_GPIOx", pinDef.port); err != nil {
				return err
			}
			if err = writeDefine(out, pin+"_GPIO_Mode", pinDef.mode); err != nil {
				return err
			}
			if err = writeDefine(out, pin+"_GPIO_PuPd", pinDef.pull); err != nil {
				return err
			}
			if err = writeDefine(out, pin+"_GPIO_Speed", pinDef.speed); err != nil {
				return err
			}
			if err = writeDefine(out, pin+"_GPIO_AF", pinDef.alternate); err != nil {
				return err
			}
		}
	}

	return nil
}

func writeDefine(out *bufio.Writer, name string, value string) error {
	invalidChars := [...]string{"=", " ", "/", "(", ")", "[", "]", "\\", "-"}

	if len(value) == 0 {
		return nil
	}
	for _, ch := range invalidChars {
		name = strings.ReplaceAll(name, ch, "_")
		value = strings.ReplaceAll(value, ch, "_")
	}
	name = "MX_" + name
	if len(name) < 39 {
		name += strings.Repeat(" ", 39-len(name))
	}
	_, err := out.WriteString("#define " + name + value + "\n")
	return err
}
