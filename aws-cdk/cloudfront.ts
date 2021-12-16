import { Construct, Duration } from "@aws-cdk/core";
import * as cloudfront from "@aws-cdk/aws-cloudfront";
import {
  AllowedMethods,
  CachedMethods,
  OriginRequestPolicy,
  SecurityPolicyProtocol,
  ViewerProtocolPolicy,
} from "@aws-cdk/aws-cloudfront";
import * as origins from "@aws-cdk/aws-cloudfront-origins";
import { Certificate } from "@aws-cdk/aws-certificatemanager";
import {
  AaaaRecord,
  ARecord,
  IHostedZone,
  RecordTarget,
} from "@aws-cdk/aws-route53";
import { CloudFrontTarget } from "@aws-cdk/aws-route53-targets";

interface ICloudFrontConstructProps {
  originEndpoint: string;
  endpoint: string;
  hostedZone: IHostedZone;
  cdnCertId: string;
}

export class CloudFront extends Construct {
  originEndpoint: string;
  endpoint: string;
  hostedZone: IHostedZone;
  cdnCertId: string;

  constructor(scope: Construct, id: string, props: ICloudFrontConstructProps) {
    super(scope, id);

    this.originEndpoint = props.originEndpoint;
    this.endpoint = props.endpoint;
    this.hostedZone = props.hostedZone;
    this.cdnCertId = props.cdnCertId;

    this.distributionDefinition();
  }

  distributionDefinition(): void {
    const cachePolicy = new cloudfront.CachePolicy(
      this,
      "soldev-api-cache-policy",
      {
        cachePolicyName: "soldev-api-cache-policy",
        comment: "A default policy for SolDev API",
        headerBehavior:
          cloudfront.CacheHeaderBehavior.allowList("Cache-Control"),
        queryStringBehavior: cloudfront.CacheQueryStringBehavior.all(),
        minTtl: Duration.seconds(0),
        defaultTtl: Duration.minutes(1),
        maxTtl: Duration.minutes(30),
        enableAcceptEncodingGzip: true,
        enableAcceptEncodingBrotli: true,
      }
    );

    const certificate = Certificate.fromCertificateArn(
      this,
      "cdn-cert",
      this.cdnCertId
    );

    const cf = new cloudfront.Distribution(this, "api-distribution", {
      defaultBehavior: {
        origin: new origins.HttpOrigin(this.originEndpoint),
        allowedMethods: AllowedMethods.ALLOW_ALL,
        cachedMethods: CachedMethods.CACHE_GET_HEAD_OPTIONS,
        compress: true,
        viewerProtocolPolicy: ViewerProtocolPolicy.HTTPS_ONLY,
        cachePolicy,
        originRequestPolicy: OriginRequestPolicy.CORS_CUSTOM_ORIGIN,
      },
      minimumProtocolVersion: SecurityPolicyProtocol.TLS_V1_2_2021,
      enableLogging: true,
      domainNames: [this.endpoint],
      certificate,
    });

    new ARecord(this, "api-cdnarecord", {
      zone: this.hostedZone,
      recordName: "api",
      target: RecordTarget.fromAlias(new CloudFrontTarget(cf)),
    });

    new AaaaRecord(this, "api-alias-record", {
      zone: this.hostedZone,
      recordName: "api",
      target: RecordTarget.fromAlias(new CloudFrontTarget(cf)),
    });
  }
}
