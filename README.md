# WISE-4050 Module Overview

## Introduction
The WISE-4050 is a 4DI/4DO 2.4G WiFi IoT Wireless I/O Module designed for a variety of applications in the Industrial Internet of Things (IIoT). It serves as a reliable and flexible solution for data acquisition and control in various industrial environments. This document provides an overview of the WISE-4050 module and its capabilities, focusing on its integration with the Modbus TCP protocol as part of its feature set.

## Features
- **Model Name:** WISE-4050
- **Description:** The WISE-4000 series is an Ethernet-based wired or wireless IoT device, integrated with IoT data acquisition, processing, and publishing functions.
- **Connectivity:** Offers both wired Ethernet and wireless 2.4G WiFi options for versatile networking solutions.
- **Digital Input/Output:** Equipped with 4 Digital Inputs (DI) and 4 Digital Outputs (DO), enabling the module to interface with various sensors and actuators.
- **Modbus TCP Support:** Integrates with Modbus TCP protocol, allowing the module to communicate with a wide range of industrial devices over a TCP/IP network.

## Modbus TCP Integration
This package includes an integration with the Modbus TCP protocol, facilitating communication between the WISE-4050 module and other industrial devices. The integration allows for:

- **Data Acquisition:** Reading sensor data from the digital inputs via Modbus TCP.
- **Device Control:** Writing to the digital outputs to control connected devices.
- **Configurable Parameters:** Flexibility to set device address, port, and coil addresses for DI and DO.

### Example Configuration for Modbus TCP:
```json
{
  "DeviceAddress": "192.168.1.100",
  "port": "502",
  "digital_inputs": {
    "base_address": 0,
    "length": 4
  },
  "digital_outputs": {
    "base_address": 16,
    "length": 4
  }
}
```
### Example DoCommands for Modbus TCP:
```json
For changing only Coil 1 and Coil 4
{
    "coil1": true,
    "coil4": false,
}
```
```json
For changing all four Coils 
{
    "coil1": true,
    "coil2": false,
    "coil3": true,
    "coil4": false,
}
```