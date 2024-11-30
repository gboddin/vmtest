all: computer test.rom
computer:
	go build -o computer ./
test.rom:
	xa test.asm -o test.rom
clean:
	rm computer test.rom
