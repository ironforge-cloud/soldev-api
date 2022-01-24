import { Construct, RemovalPolicy } from "@aws-cdk/core";
import { AttributeType, BillingMode, Table } from "@aws-cdk/aws-dynamodb";

export class DynamoDB extends Construct {
  playlistsTable: Table;
  contentTable: Table;
  usersTable: Table;
  socialTable: Table;
  bountiesTable: Table;

  constructor(scope: Construct, id: string) {
    super(scope, id);
    this.userDB();
    this.playlistDB();
    this.contentDB();
    this.socialDB();
    this.bountiesDB();
  }

  userDB(): void {
    // User Table
    this.usersTable = new Table(this, "Users", {
      partitionKey: {
        name: "PublicKey",
        type: AttributeType.STRING,
      },
      tableName: "Users",
      removalPolicy:
        process.env.AWS_ENV === "production"
          ? RemovalPolicy.RETAIN
          : RemovalPolicy.DESTROY,
      billingMode: BillingMode.PAY_PER_REQUEST,
      pointInTimeRecovery: true,
      waitForReplicationToFinish: false,
    });
  }

  playlistDB(): void {
    // Playlists Table
    this.playlistsTable = new Table(this, "Playlists", {
      partitionKey: {
        name: "Vertical",
        type: AttributeType.STRING,
      },
      sortKey: {
        name: "ID",
        type: AttributeType.STRING,
      },
      tableName: "Playlists",
      removalPolicy:
        process.env.AWS_ENV === "production"
          ? RemovalPolicy.RETAIN
          : RemovalPolicy.DESTROY,
      billingMode: BillingMode.PAY_PER_REQUEST,
      pointInTimeRecovery: true,
      waitForReplicationToFinish: false,
    });

    // GSI - Playlists Table
    // This is used in our ETL system
    this.playlistsTable.addGlobalSecondaryIndex({
      indexName: "provider-gsi",
      partitionKey: {
        name: "Provider",
        type: AttributeType.STRING,
      },
      sortKey: {
        name: "Position",
        type: AttributeType.NUMBER,
      },
    });

    // GSI - Playlists Table
    this.playlistsTable.addGlobalSecondaryIndex({
      indexName: "vertical-gsi",
      partitionKey: {
        name: "Vertical",
        type: AttributeType.STRING,
      },
      sortKey: {
        name: "Position",
        type: AttributeType.NUMBER,
      },
    });
  }

  contentDB(): void {
    // Content Table
    this.contentTable = new Table(this, "Content", {
      partitionKey: {
        name: "PK",
        type: AttributeType.STRING,
      },
      sortKey: {
        name: "SK",
        type: AttributeType.STRING,
      },
      tableName: "Content",
      timeToLiveAttribute: "Expdate",
      removalPolicy:
        process.env.AWS_ENV === "production"
          ? RemovalPolicy.RETAIN
          : RemovalPolicy.DESTROY,
      billingMode: BillingMode.PAY_PER_REQUEST,
      pointInTimeRecovery: true,
      waitForReplicationToFinish: false,
    });

    this.contentTable.addLocalSecondaryIndex({
      indexName: "content-status-lsi",
      sortKey: {
        name: "ContentStatus",
        type: AttributeType.STRING,
      },
    });

    this.contentTable.addGlobalSecondaryIndex({
      indexName: "status-gsi",
      partitionKey: {
        name: "ContentStatus",
        type: AttributeType.STRING,
      },
    });

    this.contentTable.addGlobalSecondaryIndex({
      indexName: "special-tag-gsi",
      partitionKey: {
        name: "SpecialTag",
        type: AttributeType.STRING,
      },
      sortKey: {
        name: "PublishedAt",
        type: AttributeType.STRING,
      },
    });

    this.contentTable.addGlobalSecondaryIndex({
      indexName: "video-gsi",
      sortKey: {
        name: "Position",
        type: AttributeType.NUMBER,
      },
      partitionKey: {
        name: "PK",
        type: AttributeType.STRING,
      },
    });

    this.contentTable.addGlobalSecondaryIndex({
      indexName: "url-gsi",
      partitionKey: {
        name: "Url",
        type: AttributeType.STRING,
      },
    });
  }

  socialDB(): void {
    this.socialTable = new Table(this, "Social", {
      partitionKey: {
        name: "PK",
        type: AttributeType.STRING,
      },
      sortKey: {
        name: "CreatedAt",
        type: AttributeType.STRING,
      },
      tableName: "Social",
      removalPolicy:
        process.env.AWS_ENV === "production"
          ? RemovalPolicy.RETAIN
          : RemovalPolicy.DESTROY,
      billingMode: BillingMode.PAY_PER_REQUEST,
      pointInTimeRecovery: true,
      timeToLiveAttribute: "Expdate",
      waitForReplicationToFinish: false,
    });

    this.socialTable.addGlobalSecondaryIndex({
      indexName: "tweet-gsi",
      partitionKey: {
        name: "ID",
        type: AttributeType.STRING,
      },
    });

    this.socialTable.addGlobalSecondaryIndex({
      indexName: "pinned-tweet-gsi",
      partitionKey: {
        name: "Pinned",
        type: AttributeType.NUMBER,
      },
      sortKey: {
        name: "CreatedAt",
        type: AttributeType.STRING,
      },
    });
  }

  bountiesDB(): void {
    this.bountiesTable = new Table(this, "Bounties", {
      partitionKey: {
        name: "PK",
        type: AttributeType.STRING,
      },
      sortKey: {
        name: "CreatedAt",
        type: AttributeType.STRING,
      },
      tableName: "Bounties",
      removalPolicy:
        process.env.AWS_ENV === "production"
          ? RemovalPolicy.RETAIN
          : RemovalPolicy.DESTROY,
      billingMode: BillingMode.PAY_PER_REQUEST,
      pointInTimeRecovery: true,
      waitForReplicationToFinish: false,
    });
  }
}
