export type ResponseMessage = {
  message: string;
  success: boolean;
  data: any | TableConfig;
};

export type TableConfig = {
  name: string;
  columns: ColumnConfig[];
};

export type ColumnConfig = {
  name: string;
  dataType: DataType;
  refTable: string;
  notNull: boolean;
  unique: boolean;
};

export enum DataType {
  TEXT = "TEXT",
  INTEGER = "INTEGER",
  DOUBLE = "DOUBLE",
  BOOLEAN = "BOOLEAN",
  DATETIME = "DATETIME",
  REFERENCE = "REFERENCE",
}

export type RecordModel = {
  [key: string]: any;
};

export type DataListModel = {
  collectionName: string;
  filter?: string;
  size?: number;
  limit?: number;
  ID?: number[];
};

export type DataInsertModel = {
  collectionName: string;
  values: Map<string, any>[] | RecordModel | object;
};

export type DataUpdateModel = {
  collectionName: string;
  ID: number;
  values: Map<string, any>[] | RecordModel | object;
};

export type DataDeleteModel = {
  collectionName: string;
  ID: number;
};

export type LoginModel = {
  login: string;
  password: string;
};

export type ViewModel = {
  viewName: string;
};
