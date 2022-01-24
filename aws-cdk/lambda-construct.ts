import { ITable } from "@aws-cdk/aws-dynamodb";
import { Construct, Duration } from "@aws-cdk/core";
import * as path from "path";
import { HttpApi, HttpMethod } from "@aws-cdk/aws-apigatewayv2";
import { HttpLambdaIntegration } from "@aws-cdk/aws-apigatewayv2-integrations";
import { GoFunction } from "@aws-cdk/aws-lambda-go";
import * as lambda from "@aws-cdk/aws-lambda";
import * as events from "@aws-cdk/aws-events";
import * as targets from "@aws-cdk/aws-events-targets";

interface ILambdaConstructProps {
  playlistsTable: ITable;
  contentTable: ITable;
  httpApi: HttpApi;
  usersTable: ITable;
  socialTable: ITable;
}

export class Lambda extends Construct {
  httpApi: HttpApi;
  playlistsTable: ITable;
  contentTable: ITable;
  usersTable: ITable;
  socialTable: ITable;

  constructor(scope: Construct, id: string, props: ILambdaConstructProps) {
    super(scope, id);

    this.httpApi = props.httpApi;
    this.playlistsTable = props.playlistsTable;
    this.contentTable = props.contentTable;
    this.usersTable = props.usersTable;
    this.socialTable = props.socialTable;

    // Lambda Function init
    this.GetUser();
    this.PutUser();

    this.GetPlaylists();
    this.GetPlaylistsById();
    this.PutPlaylists();

    this.PutContent();
    this.PostContent();
    this.GetContent();
    this.GetContentById();
    this.GetContentByStatus();
    this.GetContentUsingSpecialTag();
    this.GetContentUsingList();
    this.CheckContentByUrl();
    this.GetContentTypes();

    this.GetTweets();
    this.GetPinnedTweets();
    this.PinTweet();

    this.SyncYoutubeContent();
    this.SyncTwitter();

    this.ReviewNewContent();

    this.PutCompany();
    this.PostCompany();
    this.DeleteCompany();
    this.GetCompanyByID();
    this.GetAllCompanies();

    this.PostBounty();
    this.PutBounty();
    this.DeleteBounty();
    this.GetAllBountiesByCompanyID();
    this.GetBountyByID();
    this.GetBountyStatsByCompanyID();
    this.GetBountyStats();
  }

  GetBountyStats(): void {
    const lambdaFunction = new GoFunction(this, "get-bounty-stats", {
      entry: path.join(
        process.cwd(),
        "src",
        "cmd",
        "get-bounty-stats",
        "main.go"
      ),
      bundling: {
        environment: {
          GOARCH: "arm64",
          GOOS: "linux",
        },
      },
      environment: {
        POSTGRESQL_URL: process.env.POSTGRESQL_URL as string,
      },
      memorySize: 1024,
      architecture: lambda.Architecture.ARM_64,
    });

    const integration = new HttpLambdaIntegration(
      "get-bounty-stats-integration",
      lambdaFunction
    );

    this.httpApi.addRoutes({
      path: "/bounties/stats",
      methods: [HttpMethod.GET],
      integration,
    });
  }

  GetBountyStatsByCompanyID(): void {
    const lambdaFunction = new GoFunction(
      this,
      "get-bounties-stats-by-companyid",
      {
        entry: path.join(
          process.cwd(),
          "src",
          "cmd",
          "get-bounties-stats-by-companyid",
          "main.go"
        ),
        bundling: {
          environment: {
            GOARCH: "arm64",
            GOOS: "linux",
          },
        },
        environment: {
          POSTGRESQL_URL: process.env.POSTGRESQL_URL as string,
        },
        memorySize: 1024,
        architecture: lambda.Architecture.ARM_64,
      }
    );

    const integration = new HttpLambdaIntegration(
      "get-bounties-stats-by-companyid-integration",
      lambdaFunction
    );

    this.httpApi.addRoutes({
      path: "/bounties/stats/company/{companyID}",
      methods: [HttpMethod.GET],
      integration,
    });
  }

  GetBountyByID(): void {
    const lambdaFunction = new GoFunction(this, "get-bounty-by-id", {
      entry: path.join(
        process.cwd(),
        "src",
        "cmd",
        "get-bounty-by-id",
        "main.go"
      ),
      bundling: {
        environment: {
          GOARCH: "arm64",
          GOOS: "linux",
        },
      },
      environment: {
        POSTGRESQL_URL: process.env.POSTGRESQL_URL as string,
      },
      memorySize: 1024,
      architecture: lambda.Architecture.ARM_64,
    });

    const integration = new HttpLambdaIntegration(
      "get-bounty-by-id-integration",
      lambdaFunction
    );

    this.httpApi.addRoutes({
      path: "/bounties/{bountyID}",
      methods: [HttpMethod.GET],
      integration,
    });
  }

  GetAllBountiesByCompanyID(): void {
    const lambdaFunction = new GoFunction(
      this,
      "get-all-bounties-by-companyid",
      {
        entry: path.join(
          process.cwd(),
          "src",
          "cmd",
          "get-all-bounties-by-companyid",
          "main.go"
        ),
        bundling: {
          environment: {
            GOARCH: "arm64",
            GOOS: "linux",
          },
        },
        environment: {
          POSTGRESQL_URL: process.env.POSTGRESQL_URL as string,
        },
        memorySize: 1024,
        architecture: lambda.Architecture.ARM_64,
      }
    );

    const integration = new HttpLambdaIntegration(
      "get-all-bounties-by-companyid-integration",
      lambdaFunction
    );

    this.httpApi.addRoutes({
      path: "/bounties/company/{companyID}",
      methods: [HttpMethod.GET],
      integration,
    });
  }

  PutBounty(): void {
    const lambdaFunction = new GoFunction(this, "put-bounty", {
      entry: path.join(process.cwd(), "src", "cmd", "put-bounty", "main.go"),
      bundling: {
        environment: {
          GOARCH: "arm64",
          GOOS: "linux",
        },
      },
      environment: {
        POSTGRESQL_URL: process.env.POSTGRESQL_URL as string,
      },
      memorySize: 1024,
      architecture: lambda.Architecture.ARM_64,
    });

    const integration = new HttpLambdaIntegration(
      "put-bounty-integration",
      lambdaFunction
    );

    this.httpApi.addRoutes({
      path: "/bounties/{bountyID}",
      methods: [HttpMethod.PUT],
      integration,
    });
  }

  DeleteBounty(): void {
    const lambdaFunction = new GoFunction(this, "delete-bounty", {
      entry: path.join(process.cwd(), "src", "cmd", "delete-bounty", "main.go"),
      bundling: {
        environment: {
          GOARCH: "arm64",
          GOOS: "linux",
        },
      },
      environment: {
        POSTGRESQL_URL: process.env.POSTGRESQL_URL as string,
      },
      memorySize: 1024,
      architecture: lambda.Architecture.ARM_64,
    });

    const integration = new HttpLambdaIntegration(
      "delete-bounty-integration",
      lambdaFunction
    );

    this.httpApi.addRoutes({
      path: "/bounties/{bountyID}",
      methods: [HttpMethod.DELETE],
      integration,
    });
  }

  PostBounty(): void {
    const lambdaFunction = new GoFunction(this, "post-bounty", {
      entry: path.join(process.cwd(), "src", "cmd", "post-bounty", "main.go"),
      bundling: {
        environment: {
          GOARCH: "arm64",
          GOOS: "linux",
        },
      },
      environment: {
        POSTGRESQL_URL: process.env.POSTGRESQL_URL as string,
      },
      memorySize: 1024,
      architecture: lambda.Architecture.ARM_64,
    });

    const integration = new HttpLambdaIntegration(
      "post-bounty-integration",
      lambdaFunction
    );

    this.httpApi.addRoutes({
      path: "/bounties/company/{companyID}",
      methods: [HttpMethod.POST],
      integration,
    });
  }

  GetAllCompanies(): void {
    const lambdaFunction = new GoFunction(this, "get-all-companies", {
      entry: path.join(
        process.cwd(),
        "src",
        "cmd",
        "get-all-companies",
        "main.go"
      ),
      bundling: {
        environment: {
          GOARCH: "arm64",
          GOOS: "linux",
        },
      },
      environment: {
        POSTGRESQL_URL: process.env.POSTGRESQL_URL as string,
      },
      memorySize: 1024,
      architecture: lambda.Architecture.ARM_64,
    });

    const integration = new HttpLambdaIntegration(
      "get-all-companies-integration",
      lambdaFunction
    );

    this.httpApi.addRoutes({
      path: "/companies",
      methods: [HttpMethod.GET],
      integration,
    });
  }

  GetCompanyByID(): void {
    const lambdaFunction = new GoFunction(this, "get-company-by-id", {
      entry: path.join(
        process.cwd(),
        "src",
        "cmd",
        "get-company-by-id",
        "main.go"
      ),
      bundling: {
        environment: {
          GOARCH: "arm64",
          GOOS: "linux",
        },
      },
      environment: {
        POSTGRESQL_URL: process.env.POSTGRESQL_URL as string,
      },
      memorySize: 1024,
      architecture: lambda.Architecture.ARM_64,
    });

    const integration = new HttpLambdaIntegration(
      "get-company-by-id-integration",
      lambdaFunction
    );

    this.httpApi.addRoutes({
      path: "/companies/{companyID}",
      methods: [HttpMethod.GET],
      integration,
    });
  }

  DeleteCompany(): void {
    const lambdaFunction = new GoFunction(this, "delete-company", {
      entry: path.join(
        process.cwd(),
        "src",
        "cmd",
        "delete-company",
        "main.go"
      ),
      bundling: {
        environment: {
          GOARCH: "arm64",
          GOOS: "linux",
        },
      },
      environment: {
        POSTGRESQL_URL: process.env.POSTGRESQL_URL as string,
      },
      memorySize: 1024,
      architecture: lambda.Architecture.ARM_64,
    });

    const integration = new HttpLambdaIntegration(
      "delete-company-integration",
      lambdaFunction
    );

    this.httpApi.addRoutes({
      path: "/companies/{companyID}",
      methods: [HttpMethod.DELETE],
      integration,
    });
  }

  PutCompany(): void {
    const lambdaFunction = new GoFunction(this, "put-company", {
      entry: path.join(process.cwd(), "src", "cmd", "put-company", "main.go"),
      bundling: {
        environment: {
          GOARCH: "arm64",
          GOOS: "linux",
        },
      },
      environment: {
        POSTGRESQL_URL: process.env.POSTGRESQL_URL as string,
      },
      memorySize: 1024,
      architecture: lambda.Architecture.ARM_64,
    });

    const integration = new HttpLambdaIntegration(
      "put-company-integration",
      lambdaFunction
    );

    this.httpApi.addRoutes({
      path: "/companies/{companyID}",
      methods: [HttpMethod.PUT],
      integration,
    });
  }

  PostCompany(): void {
    const lambdaFunction = new GoFunction(this, "post-company", {
      entry: path.join(process.cwd(), "src", "cmd", "post-company", "main.go"),
      bundling: {
        environment: {
          GOARCH: "arm64",
          GOOS: "linux",
        },
      },
      environment: {
        POSTGRESQL_URL: process.env.POSTGRESQL_URL as string,
      },
      memorySize: 1024,
      architecture: lambda.Architecture.ARM_64,
    });

    const integration = new HttpLambdaIntegration(
      "post-company-integration",
      lambdaFunction
    );

    this.httpApi.addRoutes({
      path: "/companies",
      methods: [HttpMethod.POST],
      integration,
    });
  }

  GetUser(): void {
    const lambdaFunction = new GoFunction(this, "get-user", {
      entry: path.join(process.cwd(), "src", "cmd", "get-user", "main.go"),
      bundling: {
        environment: {
          GOARCH: "arm64",
          GOOS: "linux",
        },
      },
      memorySize: 1024,
      architecture: lambda.Architecture.ARM_64,
    });

    this.usersTable.grantReadWriteData(lambdaFunction);

    const integration = new HttpLambdaIntegration(
      "get-user-integration",
      lambdaFunction
    );

    this.httpApi.addRoutes({
      path: "/user/{publicKey}",
      methods: [HttpMethod.GET],
      integration,
    });
  }

  PutUser(): void {
    const lambdaFunction = new GoFunction(this, "put-user", {
      entry: path.join(process.cwd(), "src", "cmd", "put-user", "main.go"),
      bundling: {
        environment: {
          GOARCH: "arm64",
          GOOS: "linux",
        },
      },
      memorySize: 1024,
      architecture: lambda.Architecture.ARM_64,
    });

    this.usersTable.grantWriteData(lambdaFunction);

    const integration = new HttpLambdaIntegration(
      "put-user-integration",
      lambdaFunction
    );

    this.httpApi.addRoutes({
      path: "/user",
      methods: [HttpMethod.PUT],
      integration,
    });
  }

  GetPlaylists(): void {
    const lambdaFunction = new GoFunction(this, "get-playlists", {
      entry: path.join(process.cwd(), "src", "cmd", "get-playlists", "main.go"),
      bundling: {
        environment: {
          GOARCH: "arm64",
          GOOS: "linux",
        },
      },
      memorySize: 1024,
      architecture: lambda.Architecture.ARM_64,
    });

    this.playlistsTable.grantReadData(lambdaFunction);

    const integration = new HttpLambdaIntegration(
      "get-playlists-integration",
      lambdaFunction
    );

    this.httpApi.addRoutes({
      path: "/playlists/{vertical}",
      methods: [HttpMethod.GET],
      integration,
    });
  }

  GetPlaylistsById(): void {
    const lambdaFunction = new GoFunction(this, "get-playlists-by-id", {
      entry: path.join(
        process.cwd(),
        "src",
        "cmd",
        "get-playlists-by-id",
        "main.go"
      ),
      bundling: {
        environment: {
          GOARCH: "arm64",
          GOOS: "linux",
        },
      },
      memorySize: 1024,
      architecture: lambda.Architecture.ARM_64,
    });

    this.playlistsTable.grantReadData(lambdaFunction);

    const integration = new HttpLambdaIntegration(
      "get-playlists-by-id-integration",
      lambdaFunction
    );

    this.httpApi.addRoutes({
      path: "/playlists/{vertical}/{playlistID}",
      methods: [HttpMethod.GET],
      integration,
    });
  }

  GetContent(): void {
    const lambdaFunction = new GoFunction(this, "get-content", {
      entry: path.join(process.cwd(), "src", "cmd", "get-content", "main.go"),
      bundling: {
        environment: {
          GOARCH: "arm64",
          GOOS: "linux",
        },
      },
      memorySize: 1024,
      architecture: lambda.Architecture.ARM_64,
    });

    this.contentTable.grantReadData(lambdaFunction);

    const integration = new HttpLambdaIntegration(
      "get-content-integration",
      lambdaFunction
    );

    this.httpApi.addRoutes({
      path: "/content/{vertical}/{contentType}",
      methods: [HttpMethod.GET],
      integration,
    });
  }

  GetContentUsingSpecialTag(): void {
    const lambdaFunction = new GoFunction(
      this,
      "get-content-using-specialtag",
      {
        entry: path.join(
          process.cwd(),
          "src",
          "cmd",
          "get-content-using-specialtag",
          "main.go"
        ),
        bundling: {
          environment: {
            GOARCH: "arm64",
            GOOS: "linux",
          },
        },
        memorySize: 1024,
        architecture: lambda.Architecture.ARM_64,
      }
    );

    this.contentTable.grantReadData(lambdaFunction);

    const integration = new HttpLambdaIntegration(
      "get-content-using-special-tag-integration",
      lambdaFunction
    );

    this.httpApi.addRoutes({
      path: "/content/specialtag/{specialTag}",
      methods: [HttpMethod.GET],
      integration,
    });
  }

  GetContentUsingList(): void {
    const lambdaFunction = new GoFunction(this, "get-content-using-list", {
      entry: path.join(
        process.cwd(),
        "src",
        "cmd",
        "get-content-using-list",
        "main.go"
      ),
      bundling: {
        environment: {
          GOARCH: "arm64",
          GOOS: "linux",
        },
      },
      memorySize: 1024,
      architecture: lambda.Architecture.ARM_64,
    });

    this.contentTable.grantReadData(lambdaFunction);

    const integration = new HttpLambdaIntegration(
      "get-content-using-list-integration",
      lambdaFunction
    );

    this.httpApi.addRoutes({
      path: "/content/lists/{listName}",
      methods: [HttpMethod.GET],
      integration,
    });
  }

  GetContentById(): void {
    const lambdaFunction = new GoFunction(this, "get-content-by-id", {
      entry: path.join(
        process.cwd(),
        "src",
        "cmd",
        "get-content-by-id",
        "main.go"
      ),
      bundling: {
        environment: {
          GOARCH: "arm64",
          GOOS: "linux",
        },
      },
      memorySize: 1024,
      architecture: lambda.Architecture.ARM_64,
    });

    this.contentTable.grantReadData(lambdaFunction);

    const integration = new HttpLambdaIntegration(
      "get-content-by-id-integration",
      lambdaFunction
    );

    this.httpApi.addRoutes({
      path: "/content/{vertical}/{contentType}/{ID}",
      methods: [HttpMethod.GET],
      integration,
    });
  }

  GetContentTypes(): void {
    const lambdaFunction = new GoFunction(this, "get-content-types", {
      entry: path.join(
        process.cwd(),
        "src",
        "cmd",
        "get-content-types",
        "main.go"
      ),
      bundling: {
        environment: {
          GOARCH: "arm64",
          GOOS: "linux",
        },
      },
      memorySize: 1024,
      architecture: lambda.Architecture.ARM_64,
    });

    const integration = new HttpLambdaIntegration(
      "get-content-by-id-integration",
      lambdaFunction
    );

    this.httpApi.addRoutes({
      path: "/content/types",
      methods: [HttpMethod.GET],
      integration,
    });
  }

  CheckContentByUrl(): void {
    const lambdaFunction = new GoFunction(this, "get-check-content-by-url", {
      entry: path.join(
        process.cwd(),
        "src",
        "cmd",
        "get-check-content-by-url",
        "main.go"
      ),
      bundling: {
        environment: {
          GOARCH: "arm64",
          GOOS: "linux",
        },
      },
      memorySize: 1024,
      architecture: lambda.Architecture.ARM_64,
    });

    this.contentTable.grantReadData(lambdaFunction);

    const integration = new HttpLambdaIntegration(
      "get-content-by-url-integration",
      lambdaFunction
    );

    this.httpApi.addRoutes({
      path: "/content/check",
      methods: [HttpMethod.GET],
      integration,
    });
  }

  PutPlaylists(): void {
    const lambdaFunction = new GoFunction(this, "put-playlists", {
      entry: path.join(process.cwd(), "src", "cmd", "put-playlists", "main.go"),
      bundling: {
        environment: {
          GOARCH: "arm64",
          GOOS: "linux",
        },
      },
      memorySize: 1024,
      architecture: lambda.Architecture.ARM_64,
    });

    this.playlistsTable.grantWriteData(lambdaFunction);

    const integration = new HttpLambdaIntegration(
      "put-playlists-integration",
      lambdaFunction
    );

    this.httpApi.addRoutes({
      path: "/playlists",
      methods: [HttpMethod.PUT],
      integration,
    });
  }

  PutContent(): void {
    const lambdaFunction = new GoFunction(this, "put-content", {
      entry: path.join(process.cwd(), "src", "cmd", "put-content", "main.go"),
      bundling: {
        environment: {
          GOARCH: "arm64",
          GOOS: "linux",
        },
      },
      environment: {
        ALGOLIA_APP_ID: process.env.ALGOLIA_APP_ID as string,
        ALGOLIA_API_KEY: process.env.ALGOLIA_API_KEY as string,
        DEV_ALGOLIA_INDEX: process.env.DEV_ALGOLIA_INDEX as string,
        PROD_ALGOLIA_INDEX: process.env.PROD_ALGOLIA_INDEX as string,
        AWS_ENV: process.env.AWS_ENV as string,
      },
      timeout: Duration.seconds(30),
      memorySize: 1024,
      architecture: lambda.Architecture.ARM_64,
    });

    this.contentTable.grantWriteData(lambdaFunction);

    const integration = new HttpLambdaIntegration(
      "put-content-integration",
      lambdaFunction
    );

    this.httpApi.addRoutes({
      path: "/content",
      methods: [HttpMethod.PUT],
      integration,
    });
  }

  PostContent(): void {
    const lambdaFunction = new GoFunction(this, "post-content", {
      entry: path.join(process.cwd(), "src", "cmd", "post-content", "main.go"),
      bundling: {
        environment: {
          GOARCH: "arm64",
          GOOS: "linux",
        },
      },
      environment: {
        ALGOLIA_APP_ID: process.env.ALGOLIA_APP_ID as string,
        ALGOLIA_API_KEY: process.env.ALGOLIA_API_KEY as string,
        DEV_ALGOLIA_INDEX: process.env.DEV_ALGOLIA_INDEX as string,
        PROD_ALGOLIA_INDEX: process.env.PROD_ALGOLIA_INDEX as string,
        AWS_ENV: process.env.AWS_ENV as string,
      },
      timeout: Duration.seconds(30),
      memorySize: 1024,
      architecture: lambda.Architecture.ARM_64,
    });

    this.contentTable.grantWriteData(lambdaFunction);

    const integration = new HttpLambdaIntegration(
      "post-content-integration",
      lambdaFunction
    );

    this.httpApi.addRoutes({
      path: "/content",
      methods: [HttpMethod.POST],
      integration,
    });
  }

  SyncYoutubeContent(): void {
    const lambdaFunction = new GoFunction(this, "sync-youtube-content", {
      entry: path.join(
        process.cwd(),
        "src",
        "cmd",
        "sync-youtube-content",
        "main.go"
      ),
      environment: {
        YOUTUBE_API_KEY: process.env.YOUTUBE_API_KEY as string,
        ALGOLIA_APP_ID: process.env.ALGOLIA_APP_ID as string,
        ALGOLIA_API_KEY: process.env.ALGOLIA_API_KEY as string,
        DEV_ALGOLIA_INDEX: process.env.DEV_ALGOLIA_INDEX as string,
        PROD_ALGOLIA_INDEX: process.env.PROD_ALGOLIA_INDEX as string,
        AWS_ENV: process.env.AWS_ENV as string,
      },
      timeout: Duration.minutes(5),
      bundling: {
        environment: {
          GOARCH: "arm64",
          GOOS: "linux",
        },
      },
      memorySize: 1024,
      architecture: lambda.Architecture.ARM_64,
    });

    this.playlistsTable.grantReadData(lambdaFunction);
    this.contentTable.grantReadWriteData(lambdaFunction);

    const integration = new HttpLambdaIntegration(
      "sync-youtube-content-integration",
      lambdaFunction
    );

    this.httpApi.addRoutes({
      path: "/integrations/youtube",
      methods: [HttpMethod.GET],
      integration,
    });

    const rule = new events.Rule(this, "YoutubeCron", {
      schedule: events.Schedule.expression("rate(12 hours)"),
    });

    rule.addTarget(new targets.LambdaFunction(lambdaFunction));
  }

  GetContentByStatus(): void {
    const lambdaFunction = new GoFunction(this, "get-content-by-status", {
      entry: path.join(
        process.cwd(),
        "src",
        "cmd",
        "get-content-by-status",
        "main.go"
      ),
      bundling: {
        environment: {
          GOARCH: "arm64",
          GOOS: "linux",
        },
      },
      memorySize: 1024,
      architecture: lambda.Architecture.ARM_64,
    });

    this.contentTable.grantReadData(lambdaFunction);

    const integration = new HttpLambdaIntegration(
      "get-content-by-status-integration",
      lambdaFunction
    );

    this.httpApi.addRoutes({
      path: "/content/{status}",
      methods: [HttpMethod.GET],
      integration,
    });
  }

  SyncTwitter(): void {
    const lambdaFunction = new GoFunction(this, "sync-twitter", {
      entry: path.join(process.cwd(), "src", "cmd", "sync-twitter", "main.go"),
      environment: {
        TWITTER_BEARER_TOKEN: process.env.TWITTER_BEARER_TOKEN as string,
      },
      timeout: Duration.minutes(5),
      bundling: {
        environment: {
          GOARCH: "arm64",
          GOOS: "linux",
        },
      },
      memorySize: 1024,
      architecture: lambda.Architecture.ARM_64,
    });

    this.socialTable.grantReadWriteData(lambdaFunction);

    const integration = new HttpLambdaIntegration(
      "sync-twitter-integration",
      lambdaFunction
    );

    this.httpApi.addRoutes({
      path: "/integrations/twitter",
      methods: [HttpMethod.GET],
      integration,
    });

    const rule = new events.Rule(this, "TwitterCron", {
      schedule: events.Schedule.expression("rate(1 minute)"),
    });

    rule.addTarget(new targets.LambdaFunction(lambdaFunction));
  }

  GetTweets(): void {
    const lambdaFunction = new GoFunction(this, "get-tweets", {
      entry: path.join(process.cwd(), "src", "cmd", "get-tweets", "main.go"),
      bundling: {
        environment: {
          GOARCH: "arm64",
          GOOS: "linux",
        },
      },
      memorySize: 1024,
      architecture: lambda.Architecture.ARM_64,
    });

    this.socialTable.grantReadData(lambdaFunction);

    const integration = new HttpLambdaIntegration(
      "get-tweets-integration",
      lambdaFunction
    );

    this.httpApi.addRoutes({
      path: "/tweets/{listID}",
      methods: [HttpMethod.GET],
      integration,
    });
  }

  GetPinnedTweets(): void {
    const lambdaFunction = new GoFunction(this, "get-pinned-tweets", {
      entry: path.join(
        process.cwd(),
        "src",
        "cmd",
        "get-pinned-tweets",
        "main.go"
      ),
      bundling: {
        environment: {
          GOARCH: "arm64",
          GOOS: "linux",
        },
      },
      memorySize: 1024,
      architecture: lambda.Architecture.ARM_64,
    });

    this.socialTable.grantReadData(lambdaFunction);

    const integration = new HttpLambdaIntegration(
      "get-pinned-tweets-integration",
      lambdaFunction
    );

    this.httpApi.addRoutes({
      path: "/tweets/pinned",
      methods: [HttpMethod.GET],
      integration,
    });
  }

  PinTweet(): void {
    const lambdaFunction = new GoFunction(this, "patch-pin-tweet", {
      entry: path.join(
        process.cwd(),
        "src",
        "cmd",
        "patch-pin-tweet",
        "main.go"
      ),
      bundling: {
        environment: {
          GOARCH: "arm64",
          GOOS: "linux",
        },
      },
      memorySize: 1024,
      architecture: lambda.Architecture.ARM_64,
    });

    this.socialTable.grantReadWriteData(lambdaFunction);

    const integration = new HttpLambdaIntegration(
      "get-pinned-tweets-integration",
      lambdaFunction
    );

    this.httpApi.addRoutes({
      path: "/tweets/pin/{tweetID}",
      methods: [HttpMethod.PATCH],
      integration,
    });
  }

  ReviewNewContent(): void {
    const lambdaFunction = new GoFunction(this, "review-new-content", {
      entry: path.join(
        process.cwd(),
        "src",
        "cmd",
        "review-new-content",
        "main.go"
      ),
      bundling: {
        environment: {
          GOARCH: "arm64",
          GOOS: "linux",
        },
      },
      memorySize: 1024,
      architecture: lambda.Architecture.ARM_64,
    });

    this.contentTable.grantReadWriteData(lambdaFunction);

    const rule = new events.Rule(this, "ReviewNewContentCron", {
      schedule: events.Schedule.expression("rate(1 day)"),
    });

    rule.addTarget(new targets.LambdaFunction(lambdaFunction));
  }
}
