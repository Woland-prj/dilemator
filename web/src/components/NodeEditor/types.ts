import { UUIDTypes } from "uuid";

export type Scenario = {
  ID: UUIDTypes;
  Name: string;
};

export type Node = {
  ID: UUIDTypes;
  ParentID: UUIDTypes;
  Name: string;
  Value: string;
  Scenarios: Scenario[];
};

export type Dilemma = {
  ID: UUIDTypes;
  OwnerID: UUIDTypes;
  Topic: string;
  RootNode: Node;
};
