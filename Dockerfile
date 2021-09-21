#get an image from docker 
FROM golang:1.17.1
#create a new folder called app inside the image above
RUN mkdir /app 
#add whatever is in this folder then take it to the app folder inside the image
ADD . /app
#app is the working directory
WORKDIR /app
#build an executable from the main.go file 
RUN go build -o main .
#go inside the app directory and run the main excutable file 
CMD ["/app/main"]