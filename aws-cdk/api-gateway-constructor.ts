import { Construct } from "constructs";
import { Duration } from "aws-cdk-lib";
import { CfnStage } from "aws-cdk-lib/aws-apigatewayv2";
import {
  CorsHttpMethod,
  DomainName,
  HttpApi,
} from "@aws-cdk/aws-apigatewayv2-alpha";
import { ARecord, IHostedZone, RecordTarget } from "aws-cdk-lib/aws-route53";
import { LogGroup } from "aws-cdk-lib/aws-logs";
import { PolicyStatement, Role, ServicePrincipal } from "aws-cdk-lib/aws-iam";

interface IAPIGatewayProps {
  certId: string;
  hostedZone: IHostedZone;
  domainName: DomainName;
}

export class APIGateway extends Construct {
  httpApi: HttpApi;
  certId: string;
  domainName: DomainName;
  hostedZone: IHostedZone;

  constructor(scope: Construct, id: string, props: IAPIGatewayProps) {
    super(scope, id);

    this.certId = props.certId;
    this.hostedZone = props.hostedZone;
    this.domainName = props.domainName;

    this.apiDefinition();
    this.route53Setup();
    this.setupLogs();
  }

  apiDefinition(): void {
    this.httpApi = new HttpApi(this, "api", {
      corsPreflight: {
        allowHeaders: ["Authorization", "Content-Type", "If-None-Match"],
        allowMethods: [CorsHttpMethod.ANY],
        allowOrigins: ["*"],
        maxAge: Duration.days(10),
      },
      defaultDomainMapping: {
        domainName: this.domainName,
      },
    });
  }

  route53Setup(): void {
    const dnsName = this.domainName.regionalDomainName;
    const hostedZoneId = this.domainName.regionalHostedZoneId;

    // Creating required A Record in Route53
    new ARecord(this, "api-gateway", {
      zone: this.hostedZone,
      recordName: "gateway",
      target: RecordTarget.fromAlias({
        bind() {
          return {
            dnsName,
            hostedZoneId,
          };
        },
      }),
    });
  }

  setupLogs(): void {
    const accessLogs = new LogGroup(this, "apigateway-access-logs");
    const stage = this.httpApi.defaultStage?.node.defaultChild as CfnStage;

    stage.accessLogSettings = {
      destinationArn: accessLogs.logGroupArn,
      format: JSON.stringify({
        requestId: "$context.requestId",
        userAgent: "$context.identity.userAgent",
        sourceIp: "$context.identity.sourceIp",
        requestTime: "$context.requestTime",
        requestTimeEpoch: "$context.requestTimeEpoch",
        httpMethod: "$context.httpMethod",
        path: "$context.path",
        status: "$context.status",
        protocol: "$context.protocol",
        responseLength: "$context.responseLength",
        integrationLatency: "$context.integration.latency",
        responseLatency: "$context.responseLatency",
        domainName: "$context.domainName",
        error: "$context.error.message",
      }),
    };

    const role = new Role(this, "ApiGWLogWriterRole", {
      assumedBy: new ServicePrincipal("apigateway.amazonaws.com"),
    });

    const policy = new PolicyStatement({
      actions: [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:DescribeLogGroups",
        "logs:DescribeLogStreams",
        "logs:PutLogEvents",
        "logs:GetLogEvents",
        "logs:FilterLogEvents",
      ],
      resources: ["*"],
    });

    role.addToPolicy(policy);
    accessLogs.grantWrite(role);
  }
}
