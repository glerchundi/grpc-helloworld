//go:generate closure-compiler --js ../api --js ../../../grpc/grpc-web/javascript  --js ../../../grpc/grpc-web/third_party/closure-library --js ../../../grpc/grpc-web/third_party/grpc/third_party/protobuf --entry_point=goog:proto.helloworld.GreeterClient --dependency_mode=STRICT --js_output_file js/helloworld.js
//go:generate esc -ignore \.go -o ../internal/web/web.go -pkg web .
package web
