/* USER CODE BEGIN Header */
/**
  ******************************************************************************
  * @file    usart.c
  * @brief   This file provides code for the configuration
  *          of the USART instances.
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
/* Includes ------------------------------------------------------------------*/
#include "usart.h"

/* USER CODE BEGIN 0 */

/* USER CODE END 0 */

UART_HandleTypeDef huart11;

/* USART11 init function */

void MX_USART11_UART_Init(void)
{

  /* USER CODE BEGIN USART11_Init 0 */

  /* USER CODE END USART11_Init 0 */

  /* USER CODE BEGIN USART11_Init 1 */

  /* USER CODE END USART11_Init 1 */
  huart11.Instance = USART11;
  huart11.Init.BaudRate = 115200;
  huart11.Init.WordLength = UART_WORDLENGTH_8B;
  huart11.Init.StopBits = UART_STOPBITS_1;
  huart11.Init.Parity = UART_PARITY_NONE;
  huart11.Init.Mode = UART_MODE_TX_RX;
  huart11.Init.HwFlowCtl = UART_HWCONTROL_NONE;
  huart11.Init.OverSampling = UART_OVERSAMPLING_16;
  huart11.Init.OneBitSampling = UART_ONE_BIT_SAMPLE_DISABLE;
  huart11.Init.ClockPrescaler = UART_PRESCALER_DIV1;
  huart11.AdvancedInit.AdvFeatureInit = UART_ADVFEATURE_NO_INIT;
  if (HAL_UART_Init(&huart11) != HAL_OK)
  {
    Error_Handler();
  }
  if (HAL_UARTEx_SetTxFifoThreshold(&huart11, UART_TXFIFO_THRESHOLD_1_8) != HAL_OK)
  {
    Error_Handler();
  }
  if (HAL_UARTEx_SetRxFifoThreshold(&huart11, UART_RXFIFO_THRESHOLD_1_8) != HAL_OK)
  {
    Error_Handler();
  }
  if (HAL_UARTEx_DisableFifoMode(&huart11) != HAL_OK)
  {
    Error_Handler();
  }
  /* USER CODE BEGIN USART11_Init 2 */

  /* USER CODE END USART11_Init 2 */

}

void HAL_UART_MspInit(UART_HandleTypeDef* uartHandle)
{

  GPIO_InitTypeDef GPIO_InitStruct = {0};
  RCC_PeriphCLKInitTypeDef PeriphClkInitStruct = {0};
  if(uartHandle->Instance==USART11)
  {
  /* USER CODE BEGIN USART11_MspInit 0 */

  /* USER CODE END USART11_MspInit 0 */

  /** Initializes the peripherals clock
  */
    PeriphClkInitStruct.PeriphClockSelection = RCC_PERIPHCLK_USART11;
    PeriphClkInitStruct.Usart11ClockSelection = RCC_USART11CLKSOURCE_PCLK1;
    if (HAL_RCCEx_PeriphCLKConfig(&PeriphClkInitStruct) != HAL_OK)
    {
      Error_Handler();
    }

    /* USART11 clock enable */
    __HAL_RCC_USART11_CLK_ENABLE();

    __HAL_RCC_GPIOF_CLK_ENABLE();
    /**USART11 GPIO Configuration
    PF4     ------> USART11_RX
    PF3     ------> USART11_TX
    */
    GPIO_InitStruct.Pin = GPIO_PIN_4|GPIO_PIN_3;
    GPIO_InitStruct.Mode = GPIO_MODE_AF_PP;
    GPIO_InitStruct.Pull = GPIO_NOPULL;
    GPIO_InitStruct.Speed = GPIO_SPEED_FREQ_LOW;
    GPIO_InitStruct.Alternate = GPIO_AF7_USART11;
    HAL_GPIO_Init(GPIOF, &GPIO_InitStruct);

  /* USER CODE BEGIN USART11_MspInit 1 */

  /* USER CODE END USART11_MspInit 1 */
  }
}

void HAL_UART_MspDeInit(UART_HandleTypeDef* uartHandle)
{

  if(uartHandle->Instance==USART11)
  {
  /* USER CODE BEGIN USART11_MspDeInit 0 */

  /* USER CODE END USART11_MspDeInit 0 */
    /* Peripheral clock disable */
    __HAL_RCC_USART11_CLK_DISABLE();

    /**USART11 GPIO Configuration
    PF4     ------> USART11_RX
    PF3     ------> USART11_TX
    */
    HAL_GPIO_DeInit(GPIOF, GPIO_PIN_4|GPIO_PIN_3);

  /* USER CODE BEGIN USART11_MspDeInit 1 */

  /* USER CODE END USART11_MspDeInit 1 */
  }
}

/* USER CODE BEGIN 1 */

/* USER CODE END 1 */
