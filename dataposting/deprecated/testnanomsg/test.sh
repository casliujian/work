go build testpubsub.go
# url=tcp://127.0.0.1:40899
url=ipc://testnanomsg
./testpubsub server $url server & server=$! && sleep 1
./testpubsub client $url client0 & client0=$!
./testpubsub client $url client1 & client1=$!
./testpubsub client $url client2 & client2=$!
sleep 5
kill $server $client0 $client1 $client2
rm -f testnanomsg