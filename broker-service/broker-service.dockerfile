# this file will try and make a docker image, 
# copy over the content of our directory, build the new image, then spin up another image
# and copy over the compiled build from the first image and run it



# build a tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY brokerApp /app

CMD [ "/app/brokerApp" ]