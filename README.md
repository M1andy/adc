# ADC

ADC is a simple av data organizer, for further usage, please refer to [MetaTube](https://metatube-community.github.io/wiki/).

## Examples

Given a folder structure like this:

```
/jav
├── jav_download
│   └── ipx-177.mp4
│       └── any_dir
│           └── myed-831.mp4
├── jav_output
├── adc.exe
└── config.toml
```

ADC will organize the files into the following structure:

```
/jav
├── jav_download
├── jav_output
│   ├── 五日市芽依
│   │   └── MYED-831
│   │       └── MYED-831.mp4
│   └── 相沢みなみ
│       └── IPX-177
│           └── IPX-177.mp4
├── adc.exe
└── config.toml
```

Currently only support crawl from `javbus` and only regular av number like `ipx-177`.

Organized file structure is like this: `jav_output/actor_name/av_number/av_number.mp4`

## Usage

First, you need to specify the source directory and the output directory in the `config.toml` file.

Second, you need place `adc.exe` to any directory you like, `adc.exe` will search `config.toml` from the listed directory.

- /etc/adc/
- $HOME/.config/adc
- the same directory with `adc.exe`

After that, you can run the following command to start the program

```
# one-time run
./adc.exe 

# to run as a watch dog
./adc.exe -w
```