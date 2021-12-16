#!/usr/bin/env node
import { App } from "@aws-cdk/core";
import { API } from "./api-stack";

const app = new App();

new API(app, "soldev-api");
