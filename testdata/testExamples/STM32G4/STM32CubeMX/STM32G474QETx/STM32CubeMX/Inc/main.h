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
#include "stm32g4xx_hal.h"

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
#define TIM3_CH1_Pin GPIO_PIN_2
#define TIM3_CH1_GPIO_Port GPIOE
#define TIM3_CH2_Pin GPIO_PIN_3
#define TIM3_CH2_GPIO_Port GPIOE
#define CN2_ICL_Shutout_Pin GPIO_PIN_4
#define CN2_ICL_Shutout_GPIO_Port GPIOE
#define CN2_Dissipative_Brake_Pin GPIO_PIN_5
#define CN2_Dissipative_Brake_GPIO_Port GPIOE
#define TAMPER_KEY_Pin GPIO_PIN_13
#define TAMPER_KEY_GPIO_Port GPIOC
#define PC14_OSC32_IN_Pin GPIO_PIN_14
#define PC14_OSC32_IN_GPIO_Port GPIOC
#define PC15_OSC32_OUT_Pin GPIO_PIN_15
#define PC15_OSC32_OUT_GPIO_Port GPIOC
#define TIM5_CH2_Pin GPIO_PIN_7
#define TIM5_CH2_GPIO_Port GPIOF
#define TIM5_CH2F8_Pin GPIO_PIN_8
#define TIM5_CH2F8_GPIO_Port GPIOF
#define SPI2_SCK_Pin GPIO_PIN_9
#define SPI2_SCK_GPIO_Port GPIOF
#define CN4_Dissipative_Brake_Pin GPIO_PIN_10
#define CN4_Dissipative_Brake_GPIO_Port GPIOF
#define PF0_OSC_IN_Pin GPIO_PIN_0
#define PF0_OSC_IN_GPIO_Port GPIOF
#define PF1_OSC_OUT_Pin GPIO_PIN_1
#define PF1_OSC_OUT_GPIO_Port GPIOF
#define ADC12_IN6_Pin GPIO_PIN_0
#define ADC12_IN6_GPIO_Port GPIOC
#define ADC12_IN7_Pin GPIO_PIN_1
#define ADC12_IN7_GPIO_Port GPIOC
#define ADC12_IN8_Pin GPIO_PIN_2
#define ADC12_IN8_GPIO_Port GPIOC
#define ADC12_IN9_Pin GPIO_PIN_3
#define ADC12_IN9_GPIO_Port GPIOC
#define TIM2_CH1_Pin GPIO_PIN_0
#define TIM2_CH1_GPIO_Port GPIOA
#define ADC2_IN5_Pin GPIO_PIN_4
#define ADC2_IN5_GPIO_Port GPIOC
#define LED3_Pin GPIO_PIN_11
#define LED3_GPIO_Port GPIOF
#define ADC3_IN4_Pin GPIO_PIN_7
#define ADC3_IN4_GPIO_Port GPIOE
#define TIM1_CH1N_Pin GPIO_PIN_8
#define TIM1_CH1N_GPIO_Port GPIOE
#define TIM1_CH1_Pin GPIO_PIN_9
#define TIM1_CH1_GPIO_Port GPIOE
#define TIM1_CH2N_Pin GPIO_PIN_10
#define TIM1_CH2N_GPIO_Port GPIOE
#define TIM1_CH1E11_Pin GPIO_PIN_11
#define TIM1_CH1E11_GPIO_Port GPIOE
#define TIM1_CH3N_Pin GPIO_PIN_12
#define TIM1_CH3N_GPIO_Port GPIOE
#define TIM1_CH3_Pin GPIO_PIN_13
#define TIM1_CH3_GPIO_Port GPIOE
#define ADC4_IN1_Pin GPIO_PIN_14
#define ADC4_IN1_GPIO_Port GPIOE
#define TIM1_BKIN_Pin GPIO_PIN_15
#define TIM1_BKIN_GPIO_Port GPIOE
#define OPAMP4_VOUT_Pin GPIO_PIN_12
#define OPAMP4_VOUT_GPIO_Port GPIOB
#define SPI2_MISO_Pin GPIO_PIN_14
#define SPI2_MISO_GPIO_Port GPIOB
#define SPI2_MOSI_Pin GPIO_PIN_15
#define SPI2_MOSI_GPIO_Port GPIOB
#define ADC45_IN12_Pin GPIO_PIN_8
#define ADC45_IN12_GPIO_Port GPIOD
#define ADC45_IN13_Pin GPIO_PIN_9
#define ADC45_IN13_GPIO_Port GPIOD
#define ADC345_IN7_Pin GPIO_PIN_10
#define ADC345_IN7_GPIO_Port GPIOD
#define ADC345_IN9_Pin GPIO_PIN_12
#define ADC345_IN9_GPIO_Port GPIOD
#define ADC345_IN10_Pin GPIO_PIN_13
#define ADC345_IN10_GPIO_Port GPIOD
#define CN4_ICL_Shutout_Pin GPIO_PIN_15
#define CN4_ICL_Shutout_GPIO_Port GPIOD
#define TIM8_CH1_Pin GPIO_PIN_6
#define TIM8_CH1_GPIO_Port GPIOC
#define TIM8_CH2_Pin GPIO_PIN_7
#define TIM8_CH2_GPIO_Port GPIOC
#define TIM8_CH3_Pin GPIO_PIN_8
#define TIM8_CH3_GPIO_Port GPIOC
#define LCD_CS_Pin GPIO_PIN_9
#define LCD_CS_GPIO_Port GPIOC
#define USART1_TX_Pin GPIO_PIN_9
#define USART1_TX_GPIO_Port GPIOA
#define USART1_RX_Pin GPIO_PIN_10
#define USART1_RX_GPIO_Port GPIOA
#define USB_DM_Pin GPIO_PIN_11
#define USB_DM_GPIO_Port GPIOA
#define USB_DP_Pin GPIO_PIN_12
#define USB_DP_GPIO_Port GPIOA
#define JTMS_SWDIO_Pin GPIO_PIN_13
#define JTMS_SWDIO_GPIO_Port GPIOA
#define TIM5_CH1_Pin GPIO_PIN_6
#define TIM5_CH1_GPIO_Port GPIOF
#define JTCK_SWCLK_Pin GPIO_PIN_14
#define JTCK_SWCLK_GPIO_Port GPIOA
#define TIM8_CH1N_Pin GPIO_PIN_10
#define TIM8_CH1N_GPIO_Port GPIOC
#define TIM8_CH2N_Pin GPIO_PIN_11
#define TIM8_CH2N_GPIO_Port GPIOC
#define TIM8_CH3N_Pin GPIO_PIN_12
#define TIM8_CH3N_GPIO_Port GPIOC
#define I2C3_SMBA_Pin GPIO_PIN_6
#define I2C3_SMBA_GPIO_Port GPIOG
#define I2C3_SCL_Pin GPIO_PIN_7
#define I2C3_SCL_GPIO_Port GPIOG
#define I2C3_SDA_Pin GPIO_PIN_8
#define I2C3_SDA_GPIO_Port GPIOG
#define LED1_Pin GPIO_PIN_9
#define LED1_GPIO_Port GPIOG
#define TIM3_ETR_Pin GPIO_PIN_2
#define TIM3_ETR_GPIO_Port GPIOD
#define TIM2_CH2_Pin GPIO_PIN_4
#define TIM2_CH2_GPIO_Port GPIOD
#define TIM2_CH3_Pin GPIO_PIN_7
#define TIM2_CH3_GPIO_Port GPIOD
#define JTDO_SWO_Pin GPIO_PIN_3
#define JTDO_SWO_GPIO_Port GPIOB
#define JTRST_PD_CC2_Pin GPIO_PIN_4
#define JTRST_PD_CC2_GPIO_Port GPIOB
#define TIM8_BKIN_Pin GPIO_PIN_7
#define TIM8_BKIN_GPIO_Port GPIOB
#define BOOT0_FDCAN1_RX_Pin GPIO_PIN_8
#define BOOT0_FDCAN1_RX_GPIO_Port GPIOB
#define FDCAN1_TX_Pin GPIO_PIN_9
#define FDCAN1_TX_GPIO_Port GPIOB

/* USER CODE BEGIN Private defines */

/* USER CODE END Private defines */

#ifdef __cplusplus
}
#endif

#endif /* __MAIN_H */
