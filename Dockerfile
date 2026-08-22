FROM alpine:3.20
RUN apk --no-cache add ca-certificates tzdata
COPY x402fuel-linux /x402fuel
EXPOSE 8420
ENTRYPOINT ["/x402fuel"]
CMD ["serve"]