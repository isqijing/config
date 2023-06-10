[[English](README.md) | [中文](README.zh-CN.md)]
## 1. Project Introduction

`config`: A parameter configuration program that fits the tool program very well. At the same time, it is also part of the platform development technology integration. Welcome to participate in the joint construction.

## 2. Problems solved by the project

**The parameter configuration of the program supports default values, local command line input, local configuration file input, environment variable input, grpc central service input, and has priority**

### 1. The significance of multiple parameter configuration methods coexisting

- Local command line input: Manually enter the parameter name and parameter value on the command line interface. Mostly used for cmd and shell.
- Local configuration file: Manually configure parameters, only need to be configured once, no need to configure parameters for the next run.
- Environment variables: Similar to the `local configuration file` method, but if in a complex environment, environment variables may cause conflicts between multiple programs.
- Grpc central service: The tool program will have: a unified remote control interface for external use. Affinity automation.
- Default value: The author believes that the configuration of the tool program is equipped for people who understand software configuration. In order to expand the user group, the tool program should support no parameter configuration (that is, parameter default), click to use.

### 2. The priority of multiple parameter configuration methods

There can be multiple ways to configure parameters for a program, but the final configuration is only one. **The priority of parameter configuration methods** is a kind of default behavior in disguise. `Local command line input` > `Local configuration file` > `Environment variable` > `Grpc central service` (priority from large to small).

The reason for this priority order is that, taking `local command line input` > `local configuration file` as an example, it is considered that the manual input cost of command parameters is greater than the input cost of local configuration files (the reason is simple: local configuration files can be written only once, and no configuration is required for the next execution; command line parameter form , You have to write parameters every time, poor reusability, so it takes more time).

### 3. Function introduction

1. The program should record who set the final configuration except for the default value method configuration, such as key1 is configured by grpc, key2 is configured by local configuration file, key3 is configured by local command line.

2. The initial configuration of the configuration program (default, command line, local configuration, environment variable, grpc X on/off) is configured by environment variables by default: (enable default, enable command, enable local configuration, enable environment variable, turn off grpc).

3. In order to simplify, both key and value of the configuration use string type

