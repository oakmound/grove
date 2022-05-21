#!/usr/bin/env bash

go mod tidy 

cd components
for folder in *;
 do 
  if [ -d "$folder" ]; then
    cd "$folder"
    echo "Starting tidy for components/$folder"
    go mod tidy
    cd ..
  fi
done
cd ..


cd examples
for folder in *;
 do 
  if [ -d "$folder" ]; then
    cd "$folder"
    echo "Starting tidy for examples/$folder"
    go mod tidy
    cd ..
  fi
done
cd ..



read -t 7 -p "Tidy completed"