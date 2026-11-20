ARG GO_VERSION
FROM golang:${GO_VERSION} AS build

WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -o /out/chaosarmor-triage .

FROM scratch

USER 65532:65532
COPY --from=build /out/chaosarmor-triage /chaosarmor-triage
ENTRYPOINT ["/chaosarmor-triage"]
