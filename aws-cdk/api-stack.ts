import { Construct, Stack, StackProps } from "@aws-cdk/core";
import { APIGateway } from "./api-gateway-constructor";
import { DynamoDB } from "./dynamodb-constructor";
import { Lambda } from "./lambda-construct";
import { CloudFront } from "./cloudfront";
import { HostedZone, IHostedZone } from "@aws-cdk/aws-route53";
import { DomainName } from "@aws-cdk/aws-apigatewayv2";
import { Certificate } from "@aws-cdk/aws-certificatemanager";

export class API extends Stack {
  apiGatewayCustomDomain: string;
  certId: string;
  cdnCertId: string;
  cloudFrontCustomDomain: string;
  hostedZoneName: string;
  hostedZoneId: string;
  apiGatewayDomainName: DomainName;
  hostedZone: IHostedZone;

  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);

    if (process.env.AWS_ENV === "development") {
      this.apiGatewayCustomDomain = process.env
        .DEV_APIGATEWAY_CUSTOM_DOMAIN as string;
      this.cloudFrontCustomDomain = process.env
        .DEV_CLOUDFRONT_CUSTOM_DOMAIN as string;
      this.certId = process.env.DEV_CERT_ID as string;
      this.cdnCertId = process.env.DEV_CDN_CERT_ID as string;
      this.hostedZoneId = process.env.DEV_HOSTED_ZONE_ID as string;
      this.hostedZoneName = process.env.DEV_HOSTED_ZONE_NAME as string;
    } else if (process.env.AWS_ENV === "production") {
      this.apiGatewayCustomDomain = process.env
        .PROD_APIGATEWAY_CUSTOM_DOMAIN as string;
      this.cloudFrontCustomDomain = process.env
        .PROD_CLOUDFRONT_CUSTOM_DOMAIN as string;
      this.certId = process.env.PROD_CERT_ID as string;
      this.cdnCertId = process.env.PROD_CDN_CERT_ID as string;
      this.hostedZoneId = process.env.PROD_HOSTED_ZONE_ID as string;
      this.hostedZoneName = process.env.PROD_HOSTED_ZONE_NAME as string;
    }

    this.domainSetup();

    const apiGateway = new APIGateway(this, "api", {
      certId: this.certId,
      hostedZone: this.hostedZone,
      domainName: this.apiGatewayDomainName,
    });

    new CloudFront(this, "cloudfront", {
      originEndpoint: this.apiGatewayCustomDomain,
      endpoint: this.cloudFrontCustomDomain,
      hostedZone: this.hostedZone,
      cdnCertId: this.cdnCertId,
    });

    const database = new DynamoDB(this, "db");

    new Lambda(this, "fn", {
      playlistsTable: database.playlistsTable,
      contentTable: database.contentTable,
      httpApi: apiGateway.httpApi,
      usersTable: database.usersTable,
      socialTable: database.socialTable,
    });
  }

  domainSetup(): void {
    // Retrieving hosted zone data using ID and Name
    this.hostedZone = HostedZone.fromHostedZoneAttributes(this, "hosted-zone", {
      hostedZoneId: this.hostedZoneId,
      zoneName: this.hostedZoneName,
    });

    this.apiGatewayDomainName = new DomainName(this, "api-gateway-domain", {
      domainName: this.apiGatewayCustomDomain,
      certificate: Certificate.fromCertificateArn(
        this,
        "api-gateway-domain-cert",
        this.certId
      ),
    });
  }
}
