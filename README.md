# Flightlog

A terminal UI app to explore a collection of logged flights as a personal
archive.

[![asciicast](https://asciinema.org/a/Olq5iIgkFj3pBVhGqG8mZ5ID2.svg)](https://asciinema.org/a/Olq5iIgkFj3pBVhGqG8mZ5ID2)

## Usage

```
$ flightlog -d data
```

**Options**

- `-d` points to the data directory (`data` by default)

## Data

The app reads flights from files in a directory, each flight has a JSON file
which contains basic info which isn't part of the track:

```json
# data/KL605_2435bc05.json
{
    "From": "AMS",
    "To": "SFO",
    "Number": "KL605",
    "Operator": "KLM",
    "Aircraft": "Boeing 787-9 Dreamliner",
    "Registration": "PH-BHD",
    "ScheduledDeparture": "2020-03-17T14:45:00+08:00",
    "ScheduledArrival": "2020-03-17T15:40:00-07:00"
}
```

Departure and arrival times must be in local time with offsets if applicable.

For each JSON file there must also be a corresponding CSV file with the same
name excluding the extension which contains track data in the format FlightRadar24
exposes:

```csv
# data/KL605_2435bc05.csv
Timestamp,UTC,Callsign,Position,Altitude,Speed,Direction
1584436926,2020-03-17T09:22:06Z,KLM281,"52.314148,4.770241",0,0,271
1584437115,2020-03-17T09:25:15Z,KLM281,"52.314045,4.770499",0,14,303
1584437128,2020-03-17T09:25:28Z,KLM281,"52.313896,4.770908",0,15,309
```

## Accuracy

The app can only be as accurate as the track data when calculating statistics.