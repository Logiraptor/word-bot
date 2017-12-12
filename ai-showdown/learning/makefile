
run: wordbot.so
	python train.py

wordbot.so: *.go
	go build -o wordbot.so -buildmode=c-shared