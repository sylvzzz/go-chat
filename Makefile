NAME = go-chat


.PHONY: compile server tui web clean

compile:
		cd server; make; cd ../client; make

server:
		cd server && make run

tui:
		cd client && make run

web:
		cd web && npm run dev

clean:
		cd server; make clean; cd ../client; make clean; cd ../web; make clean
