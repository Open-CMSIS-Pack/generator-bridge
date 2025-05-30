/* USER CODE BEGIN Header */
/**
  ******************************************************************************
  * @file           : main.h
  * @brief          : Header for main.c file.
  *                   This file contains the common defines of the application.
  ******************************************************************************
  * @attention
  *
  * Copyright (c) 2025 STMicroelectronics.
  * All rights reserved.
  *
  * This software is licensed under terms that can be found in the LICENSE file
  * in the root directory of this software component.
  * If no LICENSE file comes with this software, it is provided AS-IS.
  *
  ******************************************************************************
  */
/* USER CODE END Header */

/* Define to prevent recursive inclusion -------------------------------------*/
#ifndef __MAIN_H
#define __MAIN_H

#ifdef __cplusplus
extern "C" {
#endif

/* Includes ------------------------------------------------------------------*/
#include "stm32h7xx_hal.h"

/* Private includes ----------------------------------------------------------*/
/* USER CODE BEGIN Includes */

/* USER CODE END Includes */

/* Exported types ------------------------------------------------------------*/
/* USER CODE BEGIN ET */

/* USER CODE END ET */

/* Exported constants --------------------------------------------------------*/
/* USER CODE BEGIN EC */

/* USER CODE END EC */

/* Exported macro ------------------------------------------------------------*/
/* USER CODE BEGIN EM */

/* USER CODE END EM */

/* Exported functions prototypes ---------------------------------------------*/
void Error_Handler(void);

/* USER CODE BEGIN EFP */

/* USER CODE END EFP */

/* Private defines -----------------------------------------------------------*/
#define D28_Pin GPIO_PIN_6
#define D28_GPIO_Port GPIOI
#define FMC_NBL3_Pin GPIO_PIN_5
#define FMC_NBL3_GPIO_Port GPIOI
#define FMC_NBL2_Pin GPIO_PIN_4
#define FMC_NBL2_GPIO_Port GPIOI
#define ULPI_D7_Pin GPIO_PIN_5
#define ULPI_D7_GPIO_Port GPIOB
#define LCD_B6_Pin GPIO_PIN_5
#define LCD_B6_GPIO_Port GPIOK
#define FMC_NE3_Pin GPIO_PIN_10
#define FMC_NE3_GPIO_Port GPIOG
#define QSPI_BK2_IO2_Pin GPIO_PIN_9
#define QSPI_BK2_IO2_GPIO_Port GPIOG
#define FMC_NWE_Pin GPIO_PIN_5
#define FMC_NWE_GPIO_Port GPIOD
#define FMC_NOE_Pin GPIO_PIN_4
#define FMC_NOE_GPIO_Port GPIOD
#define SDIO1_D2_Pin GPIO_PIN_10
#define SDIO1_D2_GPIO_Port GPIOC
#define D25_Pin GPIO_PIN_1
#define D25_GPIO_Port GPIOI
#define D24_Pin GPIO_PIN_0
#define D24_GPIO_Port GPIOI
#define D29_Pin GPIO_PIN_7
#define D29_GPIO_Port GPIOI
#define FMC_NBL1_Pin GPIO_PIN_1
#define FMC_NBL1_GPIO_Port GPIOE
#define I2C1_SCL_Pin GPIO_PIN_6
#define I2C1_SCL_GPIO_Port GPIOB
#define LCD_B5_Pin GPIO_PIN_4
#define LCD_B5_GPIO_Port GPIOK
#define RMII_TX_EN_Pin GPIO_PIN_11
#define RMII_TX_EN_GPIO_Port GPIOG
#define LCD_B3_Pin GPIO_PIN_15
#define LCD_B3_GPIO_Port GPIOJ
#define FMC_NWAIT_Pin GPIO_PIN_6
#define FMC_NWAIT_GPIO_Port GPIOD
#define FDCAN1_STBY_Pin GPIO_PIN_3
#define FDCAN1_STBY_GPIO_Port GPIOD
#define SDIO1_D3_Pin GPIO_PIN_11
#define SDIO1_D3_GPIO_Port GPIOC
#define D26_Pin GPIO_PIN_2
#define D26_GPIO_Port GPIOI
#define D23_Pin GPIO_PIN_15
#define D23_GPIO_Port GPIOH
#define D22_Pin GPIO_PIN_14
#define D22_GPIO_Port GPIOH
#define PDM1_CLK_Pin GPIO_PIN_2
#define PDM1_CLK_GPIO_Port GPIOE
#define FMC_NBL0_Pin GPIO_PIN_0
#define FMC_NBL0_GPIO_Port GPIOE
#define I2C1_SDA_Pin GPIO_PIN_7
#define I2C1_SDA_GPIO_Port GPIOB
#define LCD_B7_Pin GPIO_PIN_6
#define LCD_B7_GPIO_Port GPIOK
#define LCD_B4_Pin GPIO_PIN_3
#define LCD_B4_GPIO_Port GPIOK
#define RMII_TXD1_Pin GPIO_PIN_12
#define RMII_TXD1_GPIO_Port GPIOG
#define FMC_NE1_Pin GPIO_PIN_7
#define FMC_NE1_GPIO_Port GPIOD
#define SDIO1_CLK_Pin GPIO_PIN_12
#define SDIO1_CLK_GPIO_Port GPIOC
#define D27__IS42S32800G_DQ27_Pin GPIO_PIN_3
#define D27__IS42S32800G_DQ27_GPIO_Port GPIOI
#define SAI1_SCKA_Pin GPIO_PIN_5
#define SAI1_SCKA_GPIO_Port GPIOE
#define SAI1_FSA_Pin GPIO_PIN_4
#define SAI1_FSA_GPIO_Port GPIOE
#define SAI1_SDB_Pin GPIO_PIN_3
#define SAI1_SDB_GPIO_Port GPIOE
#define SDNCAS_Pin GPIO_PIN_15
#define SDNCAS_GPIO_Port GPIOG
#define LCD_DE_Pin GPIO_PIN_7
#define LCD_DE_GPIO_Port GPIOK
#define QSPI_BK2_IO3_Pin GPIO_PIN_14
#define QSPI_BK2_IO3_GPIO_Port GPIOG
#define RMII_TXD0_Pin GPIO_PIN_13
#define RMII_TXD0_GPIO_Port GPIOG
#define LCD_B2_Pin GPIO_PIN_14
#define LCD_B2_GPIO_Port GPIOJ
#define LCD_B0_Pin GPIO_PIN_12
#define LCD_B0_GPIO_Port GPIOJ
#define SDIO1_CMD_Pin GPIO_PIN_2
#define SDIO1_CMD_GPIO_Port GPIOD
#define D2_Pin GPIO_PIN_0
#define D2_GPIO_Port GPIOD
#define USB_FS1_ID_Pin GPIO_PIN_10
#define USB_FS1_ID_GPIO_Port GPIOA
#define VBUS_FS1_Pin GPIO_PIN_9
#define VBUS_FS1_GPIO_Port GPIOA
#define D21_Pin GPIO_PIN_13
#define D21_GPIO_Port GPIOH
#define D30_Pin GPIO_PIN_9
#define D30_GPIO_Port GPIOI
#define TAMPER_KEY_Pin GPIO_PIN_13
#define TAMPER_KEY_GPIO_Port GPIOC
#define MFX_IRQOUT_Pin GPIO_PIN_8
#define MFX_IRQOUT_GPIO_Port GPIOI
#define SAI1_SDA_Pin GPIO_PIN_6
#define SAI1_SDA_GPIO_Port GPIOE
#define LCD_B1_Pin GPIO_PIN_13
#define LCD_B1_GPIO_Port GPIOJ
#define D3_Pin GPIO_PIN_1
#define D3_GPIO_Port GPIOD
#define SDIO1_D0_Pin GPIO_PIN_8
#define SDIO1_D0_GPIO_Port GPIOC
#define SDIO1_D1_Pin GPIO_PIN_9
#define SDIO1_D1_GPIO_Port GPIOC
#define MCO_Pin GPIO_PIN_8
#define MCO_GPIO_Port GPIOA
#define USB_FS1_DP_Pin GPIO_PIN_12
#define USB_FS1_DP_GPIO_Port GPIOA
#define USB_FS1_DM_Pin GPIO_PIN_11
#define USB_FS1_DM_GPIO_Port GPIOA
#define D31_Pin GPIO_PIN_10
#define D31_GPIO_Port GPIOI
#define ULPI_DIR_Pin GPIO_PIN_11
#define ULPI_DIR_GPIO_Port GPIOI
#define SDCLK_Pin GPIO_PIN_8
#define SDCLK_GPIO_Port GPIOG
#define SAI1_MCLKA_Pin GPIO_PIN_7
#define SAI1_MCLKA_GPIO_Port GPIOG
#define A2_Pin GPIO_PIN_2
#define A2_GPIO_Port GPIOF
#define A1_Pin GPIO_PIN_1
#define A1_GPIO_Port GPIOF
#define A0_Pin GPIO_PIN_0
#define A0_GPIO_Port GPIOF
#define QSPI_BK1_NCS_Pin GPIO_PIN_6
#define QSPI_BK1_NCS_GPIO_Port GPIOG
#define LCD_HSYNC_Pin GPIO_PIN_12
#define LCD_HSYNC_GPIO_Port GPIOI
#define LCD_VSYNC_Pin GPIO_PIN_13
#define LCD_VSYNC_GPIO_Port GPIOI
#define LCD_CLK_Pin GPIO_PIN_14
#define LCD_CLK_GPIO_Port GPIOI
#define A3_Pin GPIO_PIN_3
#define A3_GPIO_Port GPIOF
#define A13_Pin GPIO_PIN_3
#define A13_GPIO_Port GPIOG
#define A12_Pin GPIO_PIN_2
#define A12_GPIO_Port GPIOG
#define LCD_G7_Pin GPIO_PIN_2
#define LCD_G7_GPIO_Port GPIOK
#define A5_Pin GPIO_PIN_5
#define A5_GPIO_Port GPIOF
#define A4_Pin GPIO_PIN_4
#define A4_GPIO_Port GPIOF
#define LCD_G5_Pin GPIO_PIN_0
#define LCD_G5_GPIO_Port GPIOK
#define LCD_G6_Pin GPIO_PIN_1
#define LCD_G6_GPIO_Port GPIOK
#define QSPI_BK1_IO3_Pin GPIO_PIN_6
#define QSPI_BK1_IO3_GPIO_Port GPIOF
#define QSPI_BK1_IO2_Pin GPIO_PIN_7
#define QSPI_BK1_IO2_GPIO_Port GPIOF
#define QSPI_BK1_IO0_Pin GPIO_PIN_8
#define QSPI_BK1_IO0_GPIO_Port GPIOF
#define LCD_G4_Pin GPIO_PIN_11
#define LCD_G4_GPIO_Port GPIOJ
#define ULPI_STP_Pin GPIO_PIN_0
#define ULPI_STP_GPIO_Port GPIOC
#define LED1_RGB_Pin GPIO_PIN_10
#define LED1_RGB_GPIO_Port GPIOF
#define QSPI_BK1_IO1_Pin GPIO_PIN_9
#define QSPI_BK1_IO1_GPIO_Port GPIOF
#define LCd_G3_Pin GPIO_PIN_10
#define LCd_G3_GPIO_Port GPIOJ
#define RMII_MDC_Pin GPIO_PIN_1
#define RMII_MDC_GPIO_Port GPIOC
#define DFSDM_CLK_Pin GPIO_PIN_2
#define DFSDM_CLK_GPIO_Port GPIOC
#define DFSM_DAT1_Pin GPIO_PIN_3
#define DFSM_DAT1_GPIO_Port GPIOC
#define LCD_G2_Pin GPIO_PIN_9
#define LCD_G2_GPIO_Port GPIOJ
#define QSPI_BK2_IO0_Pin GPIO_PIN_2
#define QSPI_BK2_IO0_GPIO_Port GPIOH
#define ETH_MDIO_Pin GPIO_PIN_2
#define ETH_MDIO_GPIO_Port GPIOA
#define RMII_REF_CLK_Pin GPIO_PIN_1
#define RMII_REF_CLK_GPIO_Port GPIOA
#define LCD_R1_Pin GPIO_PIN_0
#define LCD_R1_GPIO_Port GPIOJ
#define D7_Pin GPIO_PIN_10
#define D7_GPIO_Port GPIOE
#define LCD_G1_Pin GPIO_PIN_8
#define LCD_G1_GPIO_Port GPIOJ
#define LCD_G0_Pin GPIO_PIN_7
#define LCD_G0_GPIO_Port GPIOJ
#define LCD_R7_Pin GPIO_PIN_6
#define LCD_R7_GPIO_Port GPIOJ
#define QSPI_BK2_IO1_Pin GPIO_PIN_3
#define QSPI_BK2_IO1_GPIO_Port GPIOH
#define ULPI_NXT_Pin GPIO_PIN_4
#define ULPI_NXT_GPIO_Port GPIOH
#define SDNWE_Pin GPIO_PIN_5
#define SDNWE_GPIO_Port GPIOH
#define LCD_R0_Pin GPIO_PIN_15
#define LCD_R0_GPIO_Port GPIOI
#define LCD_R2_Pin GPIO_PIN_1
#define LCD_R2_GPIO_Port GPIOJ
#define A7_Pin GPIO_PIN_13
#define A7_GPIO_Port GPIOF
#define A8_Pin GPIO_PIN_14
#define A8_GPIO_Port GPIOF
#define D6_Pin GPIO_PIN_9
#define D6_GPIO_Port GPIOE
#define D8_Pin GPIO_PIN_11
#define D8_GPIO_Port GPIOE
#define ULPI_D3_Pin GPIO_PIN_10
#define ULPI_D3_GPIO_Port GPIOB
#define ULPI_D4_Pin GPIO_PIN_11
#define ULPI_D4_GPIO_Port GPIOB
#define D18_Pin GPIO_PIN_10
#define D18_GPIO_Port GPIOH
#define D19_Pin GPIO_PIN_11
#define D19_GPIO_Port GPIOH
#define D1_Pin GPIO_PIN_15
#define D1_GPIO_Port GPIOD
#define D0_Pin GPIO_PIN_14
#define D0_GPIO_Port GPIOD
#define LCD_BL_CTRL_Pin GPIO_PIN_6
#define LCD_BL_CTRL_GPIO_Port GPIOA
#define RMII_CRS_DV_Pin GPIO_PIN_7
#define RMII_CRS_DV_GPIO_Port GPIOA
#define QSPI_CLK_Pin GPIO_PIN_2
#define QSPI_CLK_GPIO_Port GPIOB
#define A6_Pin GPIO_PIN_12
#define A6_GPIO_Port GPIOF
#define A9_Pin GPIO_PIN_15
#define A9_GPIO_Port GPIOF
#define D9_Pin GPIO_PIN_12
#define D9_GPIO_Port GPIOE
#define D12_Pin GPIO_PIN_15
#define D12_GPIO_Port GPIOE
#define LCD_R6_Pin GPIO_PIN_5
#define LCD_R6_GPIO_Port GPIOJ
#define D17_Pin GPIO_PIN_9
#define D17_GPIO_Port GPIOH
#define D20_Pin GPIO_PIN_12
#define D20_GPIO_Port GPIOH
#define A16_Pin GPIO_PIN_11
#define A16_GPIO_Port GPIOD
#define A17_Pin GPIO_PIN_12
#define A17_GPIO_Port GPIOD
#define A18_Pin GPIO_PIN_13
#define A18_GPIO_Port GPIOD
#define Potentiometer_Pin GPIO_PIN_0
#define Potentiometer_GPIO_Port GPIOA
#define ULPI_CK_Pin GPIO_PIN_5
#define ULPI_CK_GPIO_Port GPIOA
#define RMII_RXD0_Pin GPIO_PIN_4
#define RMII_RXD0_GPIO_Port GPIOC
#define ULPI_D2_Pin GPIO_PIN_1
#define ULPI_D2_GPIO_Port GPIOB
#define LCD_R3_Pin GPIO_PIN_2
#define LCD_R3_GPIO_Port GPIOJ
#define SNDRAS_Pin GPIO_PIN_11
#define SNDRAS_GPIO_Port GPIOF
#define A10_Pin GPIO_PIN_0
#define A10_GPIO_Port GPIOG
#define D5_Pin GPIO_PIN_8
#define D5_GPIO_Port GPIOE
#define D10_Pin GPIO_PIN_13
#define D10_GPIO_Port GPIOE
#define SDNE1_Pin GPIO_PIN_6
#define SDNE1_GPIO_Port GPIOH
#define D16_Pin GPIO_PIN_8
#define D16_GPIO_Port GPIOH
#define ULPI_D5_Pin GPIO_PIN_12
#define ULPI_D5_GPIO_Port GPIOB
#define RS_232RX_Pin GPIO_PIN_15
#define RS_232RX_GPIO_Port GPIOB
#define D15_Pin GPIO_PIN_10
#define D15_GPIO_Port GPIOD
#define D14_Pin GPIO_PIN_9
#define D14_GPIO_Port GPIOD
#define ULPI_D0_Pin GPIO_PIN_3
#define ULPI_D0_GPIO_Port GPIOA
#define LED3_RGB_Pin GPIO_PIN_4
#define LED3_RGB_GPIO_Port GPIOA
#define RMII_RXD1_Pin GPIO_PIN_5
#define RMII_RXD1_GPIO_Port GPIOC
#define ULPI_D1_Pin GPIO_PIN_0
#define ULPI_D1_GPIO_Port GPIOB
#define LCD_R4_Pin GPIO_PIN_3
#define LCD_R4_GPIO_Port GPIOJ
#define LCD_R5_Pin GPIO_PIN_4
#define LCD_R5_GPIO_Port GPIOJ
#define A11_Pin GPIO_PIN_1
#define A11_GPIO_Port GPIOG
#define D4_Pin GPIO_PIN_7
#define D4_GPIO_Port GPIOE
#define D11_Pin GPIO_PIN_14
#define D11_GPIO_Port GPIOE
#define SDCKE1_Pin GPIO_PIN_7
#define SDCKE1_GPIO_Port GPIOH
#define ULPI_D6_Pin GPIO_PIN_13
#define ULPI_D6_GPIO_Port GPIOB
#define RS232_TX_Pin GPIO_PIN_14
#define RS232_TX_GPIO_Port GPIOB
#define D13_Pin GPIO_PIN_8
#define D13_GPIO_Port GPIOD

/* USER CODE BEGIN Private defines */

/* USER CODE END Private defines */

#ifdef __cplusplus
}
#endif

#endif /* __MAIN_H */
