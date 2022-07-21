#!/bin/bash
set -exu

script_dir=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
cd "${script_dir}"

echo -e "package icons\n// WARNING: generated file, do not modify\n\n// Icons contains all the openweathermap icons\nvar Icons = map[string][]byte{" > icons.go

# from https://openweathermap.org/weather-conditions#Weather-Condition-Codes-2
codes="01d 01n 02d 02n 03d 03n 04d 04n 09d 10d 11d 13d 50d"

rm -f *.png
for code in ${codes}
do
    wget https://openweathermap.org/img/w/"${code}".png
    # make sure to "go get github.com/cratonica/2goarray"
    go run github.com/cratonica/2goarray "Icon${code}" icons < "${code}.png" > "icon${code}.go"
    rm "${code}.png"
    echo -e "\t\"$code\": Icon${code}," >> icons.go
done

echo -e "}" >> icons.go
