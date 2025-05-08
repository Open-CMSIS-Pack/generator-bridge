/******************************************************************************
 * File Name   : MX_Device.h
 * Date        : 09/04/2025 09:00:56
 * Description : STM32Cube MX parameter definitions
 * Note        : This file is generated with a generator out of the
 *               STM32CubeMX project and its generated files (DO NOT EDIT!)
 ******************************************************************************/

#ifndef MX_DEVICE_H__
#define MX_DEVICE_H__

/* MX_Device.h version */
#define MX_DEVICE_VERSION                       0x01000000


/*------------------------------ I2C1           -----------------------------*/
#define MX_I2C1                                 1

/* Filter Settings */
#define MX_I2C1_ANF_ENABLE                      1
#define MX_I2C1_DNF                             0

/* Peripheral Clock Frequency */
#define MX_I2C1_PERIPH_CLOCK_FREQ               48000000

/* Pins */

/* I2C1_SCL */
#define MX_I2C1_SCL_Pin                         PG14
#define MX_I2C1_SCL_GPIO_Pin                    GPIO_PIN_14
#define MX_I2C1_SCL_GPIOx                       GPIOG
#define MX_I2C1_SCL_GPIO_Mode                   GPIO_MODE_AF_OD
#define MX_I2C1_SCL_GPIO_PuPd                   GPIO_NOPULL
#define MX_I2C1_SCL_GPIO_Speed                  GPIO_SPEED_FREQ_LOW
#define MX_I2C1_SCL_GPIO_AF                     GPIO_AF4_I2C1

/* I2C1_SDA */
#define MX_I2C1_SDA_Pin                         PG13
#define MX_I2C1_SDA_GPIO_Pin                    GPIO_PIN_13
#define MX_I2C1_SDA_GPIOx                       GPIOG
#define MX_I2C1_SDA_GPIO_Mode                   GPIO_MODE_AF_OD
#define MX_I2C1_SDA_GPIO_PuPd                   GPIO_NOPULL
#define MX_I2C1_SDA_GPIO_Speed                  GPIO_SPEED_FREQ_LOW
#define MX_I2C1_SDA_GPIO_AF                     GPIO_AF4_I2C1

/*------------------------------ SDMMC1         -----------------------------*/
#define MX_SDMMC1                               1

/* Mode */
#define MX_SDMMC1_MODE_MMC                      1

/* Peripheral Clock Frequency */
#define MX_SDMMC1_PERIPH_CLOCK_FREQ             128000000

/* Pins */

/* SDMMC1_CK */
#define MX_SDMMC1_CK_Pin                        PC12
#define MX_SDMMC1_CK_GPIO_Pin                   GPIO_PIN_12
#define MX_SDMMC1_CK_GPIOx                      GPIOC
#define MX_SDMMC1_CK_GPIO_Mode                  GPIO_MODE_AF_PP
#define MX_SDMMC1_CK_GPIO_PuPd                  GPIO_NOPULL
#define MX_SDMMC1_CK_GPIO_Speed                 GPIO_SPEED_FREQ_HIGH
#define MX_SDMMC1_CK_GPIO_AF                    GPIO_AF12_SDMMC1

/* SDMMC1_CMD */
#define MX_SDMMC1_CMD_Pin                       PD2
#define MX_SDMMC1_CMD_GPIO_Pin                  GPIO_PIN_2
#define MX_SDMMC1_CMD_GPIOx                     GPIOD
#define MX_SDMMC1_CMD_GPIO_Mode                 GPIO_MODE_AF_PP
#define MX_SDMMC1_CMD_GPIO_PuPd                 GPIO_NOPULL
#define MX_SDMMC1_CMD_GPIO_Speed                GPIO_SPEED_FREQ_HIGH
#define MX_SDMMC1_CMD_GPIO_AF                   GPIO_AF12_SDMMC1

/* SDMMC1_D0 */
#define MX_SDMMC1_D0_Pin                        PC8
#define MX_SDMMC1_D0_GPIO_Pin                   GPIO_PIN_8
#define MX_SDMMC1_D0_GPIOx                      GPIOC
#define MX_SDMMC1_D0_GPIO_Mode                  GPIO_MODE_AF_PP
#define MX_SDMMC1_D0_GPIO_PuPd                  GPIO_NOPULL
#define MX_SDMMC1_D0_GPIO_Speed                 GPIO_SPEED_FREQ_HIGH
#define MX_SDMMC1_D0_GPIO_AF                    GPIO_AF12_SDMMC1

/* SDMMC1_D1 */
#define MX_SDMMC1_D1_Pin                        PC9
#define MX_SDMMC1_D1_GPIO_Pin                   GPIO_PIN_9
#define MX_SDMMC1_D1_GPIOx                      GPIOC
#define MX_SDMMC1_D1_GPIO_Mode                  GPIO_MODE_AF_PP
#define MX_SDMMC1_D1_GPIO_PuPd                  GPIO_NOPULL
#define MX_SDMMC1_D1_GPIO_Speed                 GPIO_SPEED_FREQ_HIGH
#define MX_SDMMC1_D1_GPIO_AF                    GPIO_AF12_SDMMC1

/* SDMMC1_D2 */
#define MX_SDMMC1_D2_Pin                        PC10
#define MX_SDMMC1_D2_GPIO_Pin                   GPIO_PIN_10
#define MX_SDMMC1_D2_GPIOx                      GPIOC
#define MX_SDMMC1_D2_GPIO_Mode                  GPIO_MODE_AF_PP
#define MX_SDMMC1_D2_GPIO_PuPd                  GPIO_NOPULL
#define MX_SDMMC1_D2_GPIO_Speed                 GPIO_SPEED_FREQ_HIGH
#define MX_SDMMC1_D2_GPIO_AF                    GPIO_AF12_SDMMC1

/* SDMMC1_D3 */
#define MX_SDMMC1_D3_Pin                        PC11
#define MX_SDMMC1_D3_GPIO_Pin                   GPIO_PIN_11
#define MX_SDMMC1_D3_GPIOx                      GPIOC
#define MX_SDMMC1_D3_GPIO_Mode                  GPIO_MODE_AF_PP
#define MX_SDMMC1_D3_GPIO_PuPd                  GPIO_NOPULL
#define MX_SDMMC1_D3_GPIO_Speed                 GPIO_SPEED_FREQ_HIGH
#define MX_SDMMC1_D3_GPIO_AF                    GPIO_AF12_SDMMC1

/* SDMMC1_D4 */
#define MX_SDMMC1_D4_Pin                        PB8
#define MX_SDMMC1_D4_GPIO_Pin                   GPIO_PIN_8
#define MX_SDMMC1_D4_GPIOx                      GPIOB
#define MX_SDMMC1_D4_GPIO_Mode                  GPIO_MODE_AF_PP
#define MX_SDMMC1_D4_GPIO_PuPd                  GPIO_NOPULL
#define MX_SDMMC1_D4_GPIO_Speed                 GPIO_SPEED_FREQ_HIGH
#define MX_SDMMC1_D4_GPIO_AF                    GPIO_AF12_SDMMC1

/* SDMMC1_D5 */
#define MX_SDMMC1_D5_Pin                        PC0
#define MX_SDMMC1_D5_GPIO_Pin                   GPIO_PIN_0
#define MX_SDMMC1_D5_GPIOx                      GPIOC
#define MX_SDMMC1_D5_GPIO_Mode                  GPIO_MODE_AF_PP
#define MX_SDMMC1_D5_GPIO_PuPd                  GPIO_NOPULL
#define MX_SDMMC1_D5_GPIO_Speed                 GPIO_SPEED_FREQ_HIGH
#define MX_SDMMC1_D5_GPIO_AF                    GPIO_AF12_SDMMC1

/* SDMMC1_D6 */
#define MX_SDMMC1_D6_Pin                        PC6
#define MX_SDMMC1_D6_GPIO_Pin                   GPIO_PIN_6
#define MX_SDMMC1_D6_GPIOx                      GPIOC
#define MX_SDMMC1_D6_GPIO_Mode                  GPIO_MODE_AF_PP
#define MX_SDMMC1_D6_GPIO_PuPd                  GPIO_NOPULL
#define MX_SDMMC1_D6_GPIO_Speed                 GPIO_SPEED_FREQ_HIGH
#define MX_SDMMC1_D6_GPIO_AF                    GPIO_AF12_SDMMC1

/* SDMMC1_D7 */
#define MX_SDMMC1_D7_Pin                        PC7
#define MX_SDMMC1_D7_GPIO_Pin                   GPIO_PIN_7
#define MX_SDMMC1_D7_GPIOx                      GPIOC
#define MX_SDMMC1_D7_GPIO_Mode                  GPIO_MODE_AF_PP
#define MX_SDMMC1_D7_GPIO_PuPd                  GPIO_NOPULL
#define MX_SDMMC1_D7_GPIO_Speed                 GPIO_SPEED_FREQ_HIGH
#define MX_SDMMC1_D7_GPIO_AF                    GPIO_AF12_SDMMC1

/*------------------------------ SPI1           -----------------------------*/
#define MX_SPI1                                 1

/* Peripheral Clock Frequency */
#define MX_SPI1_PERIPH_CLOCK_FREQ               48000000

/* Pins */

/* SPI1_MISO */
#define MX_SPI1_MISO_Pin                        PA6
#define MX_SPI1_MISO_GPIO_Pin                   GPIO_PIN_6
#define MX_SPI1_MISO_GPIOx                      GPIOA
#define MX_SPI1_MISO_GPIO_Mode                  GPIO_MODE_AF_PP
#define MX_SPI1_MISO_GPIO_PuPd                  GPIO_NOPULL
#define MX_SPI1_MISO_GPIO_Speed                 GPIO_SPEED_FREQ_LOW
#define MX_SPI1_MISO_GPIO_AF                    GPIO_AF5_SPI1

/* SPI1_MOSI */
#define MX_SPI1_MOSI_Pin                        PA7
#define MX_SPI1_MOSI_GPIO_Pin                   GPIO_PIN_7
#define MX_SPI1_MOSI_GPIOx                      GPIOA
#define MX_SPI1_MOSI_GPIO_Mode                  GPIO_MODE_AF_PP
#define MX_SPI1_MOSI_GPIO_PuPd                  GPIO_NOPULL
#define MX_SPI1_MOSI_GPIO_Speed                 GPIO_SPEED_FREQ_LOW
#define MX_SPI1_MOSI_GPIO_AF                    GPIO_AF5_SPI1

/* SPI1_SCK */
#define MX_SPI1_SCK_Pin                         PA1
#define MX_SPI1_SCK_GPIO_Pin                    GPIO_PIN_1
#define MX_SPI1_SCK_GPIOx                       GPIOA
#define MX_SPI1_SCK_GPIO_Mode                   GPIO_MODE_AF_PP
#define MX_SPI1_SCK_GPIO_PuPd                   GPIO_NOPULL
#define MX_SPI1_SCK_GPIO_Speed                  GPIO_SPEED_FREQ_LOW
#define MX_SPI1_SCK_GPIO_AF                     GPIO_AF5_SPI1

/*------------------------------ USART2         -----------------------------*/
#define MX_USART2                               1

/* Virtual mode */
#define MX_USART2_VM                            VM_ASYNC
#define MX_USART2_VM_ASYNC                      1

/* Pins */

/* USART2_RX */
#define MX_USART2_RX_Pin                        PA3
#define MX_USART2_RX_GPIO_Pin                   GPIO_PIN_3
#define MX_USART2_RX_GPIOx                      GPIOA
#define MX_USART2_RX_GPIO_Mode                  GPIO_MODE_AF_PP
#define MX_USART2_RX_GPIO_PuPd                  GPIO_NOPULL
#define MX_USART2_RX_GPIO_Speed                 GPIO_SPEED_FREQ_LOW
#define MX_USART2_RX_GPIO_AF                    GPIO_AF7_USART2

/* USART2_TX */
#define MX_USART2_TX_Pin                        PA2
#define MX_USART2_TX_GPIO_Pin                   GPIO_PIN_2
#define MX_USART2_TX_GPIOx                      GPIOA
#define MX_USART2_TX_GPIO_Mode                  GPIO_MODE_AF_PP
#define MX_USART2_TX_GPIO_PuPd                  GPIO_NOPULL
#define MX_USART2_TX_GPIO_Speed                 GPIO_SPEED_FREQ_LOW
#define MX_USART2_TX_GPIO_AF                    GPIO_AF7_USART2

/*------------------------------ USB_OTG_HS     -----------------------------*/
#define MX_USB_OTG_HS                           1

/* Handle */
#define MX_USB_OTG_HS_HANDLE                    hpcd_USB_OTG_HS

/* Virtual mode */
#define MX_USB_OTG_HS_VM                        Device_HS
#define MX_USB_OTG_HS_Device_HS                 1


#endif  /* MX_DEVICE_H__ */
