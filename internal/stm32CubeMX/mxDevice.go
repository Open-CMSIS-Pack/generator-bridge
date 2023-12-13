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
	"path"
	"path/filepath"
	"regexp"
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

	contexts := getContexts(contextMap)

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

	projectIndex := 0
	if len(contexts) == 0 {
		msp := path.Join(workDirAbs, mspFolder, mspName)
		fMsp, err := os.Open(msp)
		if err != nil {
			return err
		}
		defer fMsp.Close()

		subsystem := &params.Subsystem[projectIndex]
		fName := "MX_Device.h"
		fPath := path.Join(path.Dir(workDir), "drv_cfg", subsystem.SubsystemIdx.Project)
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

		err = mxDeviceWriteHeader(fMxDevice, fName)
		if err != nil {
			return err
		}
		peripherals, err := getPeripherals1(contextMap)
		if err != nil {
			return err
		}
		for _, peripheral := range peripherals {
			vmode := getVirtualMode(contextMap, peripheral)
			pins, err := getPins(contextMap, fMsp, peripheral)
			if err != nil {
				return err
			}
			err = mxDeviceWritePeripheralCfg(fMxDevice, peripheral, vmode, pins)
			if err != nil {
				return err
			}
		}
		_, err = fMxDevice.WriteString("\n#endif  /* __MX_DEVICE_H */\n")
		if err != nil {
			return err
		}
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

			subsystem := &params.Subsystem[projectIndex]
			fName := "MX_Device.h"
			fPath := path.Join(path.Dir(workDir), "drv_cfg", subsystem.SubsystemIdx.Project)
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

			projectIndex += 1

			err = mxDeviceWriteHeader(fMxDevice, fName)
			if err != nil {
				_ = fMxDevice.Close()
				return err
			}
			peripherals, err := getPeripherals(contextMap, context)
			if err != nil {
				_ = fMxDevice.Close()
				return err
			}
			for _, peripheral := range peripherals {
				vmode := getVirtualMode(contextMap, peripheral)
				pins, err := getPins(contextMap, fMsp, peripheral)
				if err != nil {
					_ = fMxDevice.Close()
					return err
				}
				err = mxDeviceWritePeripheralCfg(fMxDevice, peripheral, vmode, pins)
				if err != nil {
					_ = fMxDevice.Close()
					return err
				}
			}
			_, err = fMxDevice.WriteString("\n#endif  /* __MX_DEVICE_H */\n")
			if err != nil {
				return err
			}
			err = fMxDevice.Close()
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
func getContexts(contextMap map[string]map[string]string) []string {
	contexts := []string{}
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
	return contexts
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
	peripherals := []string{}
	mcu := contextMap["Mcu"]
	if mcu != nil {
		for ip, peri := range mcu {
			if strings.HasPrefix(ip, "IP") {
				for _, peripheral := range PERIPHERALS {
					if strings.HasPrefix(peri, peripheral) {
						peripherals = append(peripherals, peri)
						break
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
	var peripherals []string
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
		for ip, peri := range mcu {
			if strings.HasPrefix(ip, "IP") {
				if strings.Contains(contextIpsLine, peri) {
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
				label := signal["GPIO_Label"]
				if len(label) > 0 {
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
						val = strings.TrimRight(strings.TrimLeft(val, " "), " ")
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

func mxDeviceWriteHeader(fMxDevice *os.File, fName string) error {
	now := time.Now()
	dtString := now.Format("02/01/2006 15:04:05")

	var err error

	if _, err = fMxDevice.WriteString("/******************************************************************************\n"); err != nil {
		return err
	}
	if _, err = fMxDevice.WriteString(" * File Name   : " + fName + "\n"); err != nil {
		return err
	}
	if _, err = fMxDevice.WriteString(" * Date        : " + dtString + "\n"); err != nil {
		return err
	}
	if _, err = fMxDevice.WriteString(" * Description : STM32Cube MX parameter definitions\n"); err != nil {
		return err
	}
	if _, err = fMxDevice.WriteString(" * Note        : This file is generated with a generator out of the\n"); err != nil {
		return err
	}
	if _, err = fMxDevice.WriteString(" *               STM32CubeMX project and its generated files (DO NOT EDIT!)\n"); err != nil {
		return err
	}
	if _, err = fMxDevice.WriteString(" ******************************************************************************/\n\n"); err != nil {
		return err
	}
	if _, err = fMxDevice.WriteString("#ifndef __MX_DEVICE_H\n"); err != nil {
		return err
	}
	_, err = fMxDevice.WriteString("#define __MX_DEVICE_H\n\n")
	return err
}

func mxDeviceWritePeripheralCfg(fMxDevice *os.File, peripheral string, vmode string, pins map[string]PinDefinition) error {
	str := "\n/*------------------------------ " + peripheral
	if len(str) < 49 {
		str += strings.Repeat(" ", 49-len(str))
	}
	str += "-----------------------------*/\n"
	_, err := fMxDevice.WriteString(str)
	if err != nil {
		return err
	}
	if err = writeDefine(fMxDevice, peripheral, "1\n"); err != nil {
		return err
	}
	if vmode != "" {
		if _, err = fMxDevice.WriteString("/* Virtual mode */\n"); err != nil {
			return err
		}
		if err = writeDefine(fMxDevice, peripheral+"_VM", vmode); err != nil {
			return err
		}
		if err = writeDefine(fMxDevice, peripheral+"_"+vmode, "1"); err != nil {
			return err
		}
	}
	if len(pins) != 0 {
		_, err = fMxDevice.WriteString("/* Pins */\n")
		if err != nil {
			return err
		}
		for pin, pinDef := range pins {
			_, err = fMxDevice.WriteString("\n/* " + pin + " */\n")
			if err != nil {
				return err
			}
			if err = writeDefine(fMxDevice, pin+"_Pin", pinDef.p); err != nil {
				return err
			}
			if err = writeDefine(fMxDevice, pin+"_GPIO_Pin", pinDef.pin); err != nil {
				return err
			}
			if err = writeDefine(fMxDevice, pin+"_GPIOx", pinDef.port); err != nil {
				return err
			}
			if err = writeDefine(fMxDevice, pin+"_GPIO_Mode", pinDef.mode); err != nil {
				return err
			}
			if err = writeDefine(fMxDevice, pin+"_GPIO_PuPd", pinDef.pull); err != nil {
				return err
			}
			if err = writeDefine(fMxDevice, pin+"_GPIO_Speed", pinDef.speed); err != nil {
				return err
			}
			if err = writeDefine(fMxDevice, pin+"_GPIO_AF", pinDef.alternate); err != nil {
				return err
			}
		}
	}

	return nil
}

func writeDefine(fMxDevice *os.File, name string, value string) error {
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
	_, err := fMxDevice.WriteString("#define " + name + value + "\n")
	return err
}
