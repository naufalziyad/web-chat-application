# Start from golang
FROM golang

# Add Maintainer Info
LABEL maintainer="Naufal Ziyad Luthfiansyah <naufal.ziyad@detik.com>"

ENV GO111MODULE=on

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/github.com/naufalziyad/web-chat-application

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

# COPY go.mod . 
# COPY go.sum .

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download

# Install the package
RUN go build

# This container exposes port 8080 to the outside world
EXPOSE 8999

# Run the executable
CMD ["./web-chat-application"]
