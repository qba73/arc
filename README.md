[![Go](https://github.com/qba73/arct/actions/workflows/go.yml/badge.svg)](https://github.com/qba73/arct/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/qba73/arct)](https://goreportcard.com/report/github.com/qba73/arct)
[![Maintainability](https://api.codeclimate.com/v1/badges/219a1d97ec3a5d79ca6e/maintainability)](https://codeclimate.com/github/qba73/arct/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/219a1d97ec3a5d79ca6e/test_coverage)](https://codeclimate.com/github/qba73/arct/test_coverage)
![GitHub](https://img.shields.io/github/license/qba73/arct)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/qba73/arct)

# arct

A simple arc tool data transformer

The arct command-line utility was developed to solve a business problem related to log data processing.

The tool allows GIS and Data Analysts to extract specific data generated by a custom ESRI [ArcGIS tool](https://pro.arcgis.com/en/pro-app/latest/tool-reference/analysis/an-overview-of-the-analysis-toolbox.htm) (Arc Tool) and save the data in a CSV format. In the following step, analysts import the generated file to an EXCEL spreadsheet for further analysis.

The main requirements/assumptions are that input files are not larger than 10 MB; the analysts run the tool manually and manually import the CSV output file into spreadsheets.

## Usage

Getting help

```
$ arct -h
Usage of ./arct:
  -in string
        ArcTool log file to process
  -out string
        Output CSV file
  -version
        Show version
```
