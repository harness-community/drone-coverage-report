#!/bin/bash

# Print the JVM version installed
echo "Printing the installed JVM version:"
java -version

DWP=/harness

find $DWP -iname "*.exec*"
find $DWP -iname "*.class*"
find $DWP -iname "*.java*"

