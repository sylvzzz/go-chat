NAME = go-chat

compile:
		cd server; make; cd ../client; make

server:
		cd server && make run

client:
		cd client && make run

clean:
		cd server; rm $(NAME); cd ../client; rm $(NAME);
