FROM gcr.io/distroless/base-debian12:nonroot
ENV PORT=8080
COPY bin/todo-app /server
EXPOSE 8080
ENTRYPOINT ["/server"]
