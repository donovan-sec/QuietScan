#!/bin/bash
# Convert icon.png to icon.ico for Windows executable

echo "Converting assets/icon.png to assets/icon.ico..."

go run -mod=mod ./tools/convert-icon assets/icon.png assets/icon.ico

if [ $? -eq 0 ]; then
    echo "✓ Successfully created assets/icon.ico"
    echo "The icon is now ready for Windows builds!"
else
    echo "✗ Conversion failed"
    exit 1
fi


