import { GraphQLResolveInfo } from 'graphql';
import { ResolverContext } from './src/services/graphql/index';
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
export type MakeEmpty<T extends { [key: string]: unknown }, K extends keyof T> = { [_ in K]?: never };
export type Incremental<T> = T | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };
export type RequireFields<T, K extends keyof T> = Omit<T, K> & { [P in K]-?: NonNullable<T[P]> };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string; }
  String: { input: string; output: string; }
  Boolean: { input: boolean; output: boolean; }
  Int: { input: number; output: number; }
  Float: { input: number; output: number; }
};

export type Context = {
  __typename?: 'Context';
  context: Scalars['String']['output'];
  locale: Scalars['String']['output'];
  personality: Scalars['String']['output'];
  timezone: Scalars['String']['output'];
};

export type ContextInput = {
  context: Scalars['String']['input'];
  locale: Scalars['String']['input'];
  personality: Scalars['String']['input'];
  timezone: Scalars['String']['input'];
};

export type Message = {
  __typename?: 'Message';
  content: Scalars['String']['output'];
  timestamp: Scalars['String']['output'];
  type: MessageType;
};

export enum MessageType {
  AiMessage = 'AI_MESSAGE',
  UserMessage = 'USER_MESSAGE'
}

export type Mutation = {
  __typename?: 'Mutation';
  upsertContext: Context;
};


export type MutationUpsertContextArgs = {
  input: ContextInput;
};

export type Query = {
  __typename?: 'Query';
  context?: Maybe<Context>;
  messages: Array<Message>;
};

export enum Ws_InfoType {
  AudioStreamCompleted = 'AUDIO_STREAM_COMPLETED',
  AuthenticationFailed = 'AUTHENTICATION_FAILED',
  AuthenticationRequired = 'AUTHENTICATION_REQUIRED',
  AuthenticationSuccess = 'AUTHENTICATION_SUCCESS',
  AuthenticationTimeout = 'AUTHENTICATION_TIMEOUT',
  MessageStreamCompleted = 'MESSAGE_STREAM_COMPLETED',
  TranscriptionFailed = 'TRANSCRIPTION_FAILED'
}

export type Ws_Message = {
  __typename?: 'WS_Message';
  data: Scalars['String']['output'];
  timestamp: Scalars['String']['output'];
  type: Ws_MessageType;
};

export enum Ws_MessageType {
  Auth = 'AUTH',
  ConfigAudioFormat = 'CONFIG_AUDIO_FORMAT',
  DownloadAudio = 'DOWNLOAD_AUDIO',
  Info = 'INFO',
  TextChat = 'TEXT_CHAT',
  TextStream = 'TEXT_STREAM',
  UploadAudio = 'UPLOAD_AUDIO'
}

export type WithIndex<TObject> = TObject & Record<string, any>;
export type ResolversObject<TObject> = WithIndex<TObject>;

export type ResolverTypeWrapper<T> = Promise<T> | T;


export type ResolverWithResolve<TResult, TParent, TContext, TArgs> = {
  resolve: ResolverFn<TResult, TParent, TContext, TArgs>;
};
export type Resolver<TResult, TParent = {}, TContext = {}, TArgs = {}> = ResolverFn<TResult, TParent, TContext, TArgs> | ResolverWithResolve<TResult, TParent, TContext, TArgs>;

export type ResolverFn<TResult, TParent, TContext, TArgs> = (
  parent: TParent,
  args: TArgs,
  context: TContext,
  info: GraphQLResolveInfo
) => Promise<TResult> | TResult;

export type SubscriptionSubscribeFn<TResult, TParent, TContext, TArgs> = (
  parent: TParent,
  args: TArgs,
  context: TContext,
  info: GraphQLResolveInfo
) => AsyncIterable<TResult> | Promise<AsyncIterable<TResult>>;

export type SubscriptionResolveFn<TResult, TParent, TContext, TArgs> = (
  parent: TParent,
  args: TArgs,
  context: TContext,
  info: GraphQLResolveInfo
) => TResult | Promise<TResult>;

export interface SubscriptionSubscriberObject<TResult, TKey extends string, TParent, TContext, TArgs> {
  subscribe: SubscriptionSubscribeFn<{ [key in TKey]: TResult }, TParent, TContext, TArgs>;
  resolve?: SubscriptionResolveFn<TResult, { [key in TKey]: TResult }, TContext, TArgs>;
}

export interface SubscriptionResolverObject<TResult, TParent, TContext, TArgs> {
  subscribe: SubscriptionSubscribeFn<any, TParent, TContext, TArgs>;
  resolve: SubscriptionResolveFn<TResult, any, TContext, TArgs>;
}

export type SubscriptionObject<TResult, TKey extends string, TParent, TContext, TArgs> =
  | SubscriptionSubscriberObject<TResult, TKey, TParent, TContext, TArgs>
  | SubscriptionResolverObject<TResult, TParent, TContext, TArgs>;

export type SubscriptionResolver<TResult, TKey extends string, TParent = {}, TContext = {}, TArgs = {}> =
  | ((...args: any[]) => SubscriptionObject<TResult, TKey, TParent, TContext, TArgs>)
  | SubscriptionObject<TResult, TKey, TParent, TContext, TArgs>;

export type TypeResolveFn<TTypes, TParent = {}, TContext = {}> = (
  parent: TParent,
  context: TContext,
  info: GraphQLResolveInfo
) => Maybe<TTypes> | Promise<Maybe<TTypes>>;

export type IsTypeOfResolverFn<T = {}, TContext = {}> = (obj: T, context: TContext, info: GraphQLResolveInfo) => boolean | Promise<boolean>;

export type NextResolverFn<T> = () => Promise<T>;

export type DirectiveResolverFn<TResult = {}, TParent = {}, TContext = {}, TArgs = {}> = (
  next: NextResolverFn<TResult>,
  parent: TParent,
  args: TArgs,
  context: TContext,
  info: GraphQLResolveInfo
) => TResult | Promise<TResult>;



/** Mapping between all available schema types and the resolvers types */
export type ResolversTypes = ResolversObject<{
  Boolean: ResolverTypeWrapper<Scalars['Boolean']['output']>;
  Context: ResolverTypeWrapper<Context>;
  ContextInput: ContextInput;
  Message: ResolverTypeWrapper<Message>;
  MessageType: MessageType;
  Mutation: ResolverTypeWrapper<{}>;
  Query: ResolverTypeWrapper<{}>;
  String: ResolverTypeWrapper<Scalars['String']['output']>;
  WS_InfoType: Ws_InfoType;
  WS_Message: ResolverTypeWrapper<Ws_Message>;
  WS_MessageType: Ws_MessageType;
}>;

/** Mapping between all available schema types and the resolvers parents */
export type ResolversParentTypes = ResolversObject<{
  Boolean: Scalars['Boolean']['output'];
  Context: Context;
  ContextInput: ContextInput;
  Message: Message;
  Mutation: {};
  Query: {};
  String: Scalars['String']['output'];
  WS_Message: Ws_Message;
}>;

export type ContextResolvers<ContextType = ResolverContext, ParentType extends ResolversParentTypes['Context'] = ResolversParentTypes['Context']> = ResolversObject<{
  context?: Resolver<ResolversTypes['String'], ParentType, ContextType>;
  locale?: Resolver<ResolversTypes['String'], ParentType, ContextType>;
  personality?: Resolver<ResolversTypes['String'], ParentType, ContextType>;
  timezone?: Resolver<ResolversTypes['String'], ParentType, ContextType>;
  __isTypeOf?: IsTypeOfResolverFn<ParentType, ContextType>;
}>;

export type MessageResolvers<ContextType = ResolverContext, ParentType extends ResolversParentTypes['Message'] = ResolversParentTypes['Message']> = ResolversObject<{
  content?: Resolver<ResolversTypes['String'], ParentType, ContextType>;
  timestamp?: Resolver<ResolversTypes['String'], ParentType, ContextType>;
  type?: Resolver<ResolversTypes['MessageType'], ParentType, ContextType>;
  __isTypeOf?: IsTypeOfResolverFn<ParentType, ContextType>;
}>;

export type MutationResolvers<ContextType = ResolverContext, ParentType extends ResolversParentTypes['Mutation'] = ResolversParentTypes['Mutation']> = ResolversObject<{
  upsertContext?: Resolver<ResolversTypes['Context'], ParentType, ContextType, RequireFields<MutationUpsertContextArgs, 'input'>>;
}>;

export type QueryResolvers<ContextType = ResolverContext, ParentType extends ResolversParentTypes['Query'] = ResolversParentTypes['Query']> = ResolversObject<{
  context?: Resolver<Maybe<ResolversTypes['Context']>, ParentType, ContextType>;
  messages?: Resolver<Array<ResolversTypes['Message']>, ParentType, ContextType>;
}>;

export type Ws_MessageResolvers<ContextType = ResolverContext, ParentType extends ResolversParentTypes['WS_Message'] = ResolversParentTypes['WS_Message']> = ResolversObject<{
  data?: Resolver<ResolversTypes['String'], ParentType, ContextType>;
  timestamp?: Resolver<ResolversTypes['String'], ParentType, ContextType>;
  type?: Resolver<ResolversTypes['WS_MessageType'], ParentType, ContextType>;
  __isTypeOf?: IsTypeOfResolverFn<ParentType, ContextType>;
}>;

export type Resolvers<ContextType = ResolverContext> = ResolversObject<{
  Context?: ContextResolvers<ContextType>;
  Message?: MessageResolvers<ContextType>;
  Mutation?: MutationResolvers<ContextType>;
  Query?: QueryResolvers<ContextType>;
  WS_Message?: Ws_MessageResolvers<ContextType>;
}>;

