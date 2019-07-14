FROM golang
RUN go get -u github.com/jinzhu/gorm
RUN go get -u github.com/gorilla/mux
RUN go get -u github.com/go-sql-driver/mysql
COPY . /app/
WORKDIR /app/
CMD ["go", "run", "."]
