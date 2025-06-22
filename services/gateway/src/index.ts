import { ApolloGateway, IntrospectAndCompose, RemoteGraphQLDataSource } from '@apollo/gateway'
import { ApolloServer } from 'apollo-server'
import { config } from 'dotenv'
import { cert, initializeApp } from 'firebase-admin/app'
import { getAuth } from 'firebase-admin/auth'
import fetch from 'node-fetch'
import path, { join } from 'path'
import { getRoot } from './utils'

const env = process.env.TENK_ENV || 'development'
config({
  path: join(getRoot(), '.env.' + env)
})

const subgraphs = [
  { name: 'core', url: process.env.CORE_URL || 'http://localhost:8080/graphql' },
  { name: 'analytic', url: process.env.ANALYTIC_URL || 'http://localhost:8082/graphql' },
  { name: 'notification', url: process.env.NOTIFICATION_URL || 'http://localhost:8084/graphql' },
  { name: 'currency', url: process.env.CURRENCY_URL || 'http://localhost:8085/graphql' },
  { name: 'penk', url: process.env.PENK_URL || 'http://localhost:8099/graphql' }
]

async function startGateway() {
  // Check if all subgraphs are reachable
  const reachableSubgraphs: { name: string; url: string }[] = []
  await Promise.all(
    subgraphs.map(async ({ name, url }) => {
      try {
        const response = await fetch(url)
        if (response.ok) {
          console.log(`Subgraph at ${url} is reachable`)
        }

        reachableSubgraphs.push({ name, url })
      } catch (e) {
        console.error(`Subgraph at ${url} is not reachable`)
      }
    })
  )

  // Create the apollo gateway
  const gateway = new ApolloGateway({
    supergraphSdl: new IntrospectAndCompose({
      subgraphs: reachableSubgraphs
    }),
    buildService({ name, url }) {
      return new RemoteGraphQLDataSource({
        url,
        willSendRequest({ request, context }) {
          request.http?.headers.append('Authorization', context.authorization || '')
          request.http?.headers.append('X-Device-Id', context.deviceId || '')
          request.http?.headers.append('X-User-Id', context.userId || '')
        }
      })
    }
  })

  // Initialize Firebase app
  initializeApp({
    credential: cert(path.join(getRoot(), process.env.FIREBASE_ADMIN || ''))
  })

  // Start the server with the gateway
  const server = new ApolloServer({
    gateway,
    context: async ({ req }) => {
      const token = req.headers.authorization?.split(' ')[1]

      // Verify token with Firebase
      let decodedPayload = null
      if (token) {
        decodedPayload = await getAuth().verifyIdToken(token)
      }

      return {
        authorization: req.headers.authorization,
        deviceId: req.headers['x-device-id'],
        userId: decodedPayload?.uid
      }
    }
  })

  server.listen({ port: 8070 }).then(({ url }) => {
    console.log(`🚀 Server ready at ${url}`)
  })
}

startGateway().catch((error) => {
  console.error('Error starting gateway:', error)
  process.exit(1)
})
