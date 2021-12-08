#!/bin/bash

protoc --proto_path=proto --go_out=bep_proto/buildeventstream --go_opt=paths=source_relative proto/build_event_stream.proto
protoc --proto_path=proto --go_out=bep_proto --go_opt=paths=source_relative proto/command_line.proto proto/failure_details.proto invocation_policy.proto option_filters.proto