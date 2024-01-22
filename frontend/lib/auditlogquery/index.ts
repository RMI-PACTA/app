import { type AuditLogQueryReq, type AuditLogQuerySort, type AuditLogQueryWhere, type AuditLogQuerySortBy } from '@/openapi/generated/pacta'
import { type WritableComputedRef } from 'vue'
import { type LocalePathFunction } from 'vue-i18n-routing'
import { computed } from 'vue'

const encodeAuditLogQuerySorts = (sorts: AuditLogQuerySort[]): string => {
  const components: string[] = []
  for (const sort of sorts) {
    components.push(`${sort.ascending ? 'A' : 'D'}:${sort.by.replace('AuditLogQuerySortBy', '')}`)
  }
  return components.join(',')
}

const decodeAuditLogQuerySorts = (sorts: string): AuditLogQuerySort[] => {
  const components = sorts.split(',')
  const result: AuditLogQuerySort[] = []
  for (const component of components) {
    if (component === '') {
      continue
    }
    const [dir, byStr] = component.split(':')
    result.push({
      ascending: dir === 'A',
      by: 'AuditLogQuerySortBy' + byStr as AuditLogQuerySortBy,
    })
  }
  return result
}

const encodeAuditLogQueryWheres = (wheres: AuditLogQueryWhere[]): string => {
  const components: string[] = []
  for (const where of wheres) {
    if (where.inAction) {
      components.push(`Action:${where.inAction.join('|')}`)
    } else if (where.inActorId) {
      components.push(`ActorId:${where.inActorId.join('|')}`)
    } else if (where.inActorType) {
      components.push(`ActorType:${where.inActorType.join('|')}`)
    } else if (where.inActorOwnerId) {
      components.push(`ActorOwnerId:${where.inActorOwnerId.join('|')}`)
    } else if (where.inId) {
      components.push(`Id:${where.inId.join('|')}`)
    } else if (where.inTargetType) {
      components.push(`TargetType:${where.inTargetType.join('|')}`)
    } else if (where.inTargetId) {
      components.push(`TargetId:${where.inTargetId.join('|')}`)
    } else if (where.inTargetOwnerId) {
      components.push(`TargetOwnerId:${where.inTargetOwnerId.join('|')}`)
    } else if (where.minCreatedAt) {
      components.push(`MinCreatedAt:${where.minCreatedAt.replaceAll(':', '_')}`)
    } else if (where.maxCreatedAt) {
      components.push(`MaxCreatedAt:${where.maxCreatedAt.replaceAll(':', '_')}`)
    } else {
      console.warn(new Error(`Unknown where: ${JSON.stringify(where)}`))
    }
  }
  return components.join(',')
}

const decodeAudtLogQueryWheres = (wheres: string): AuditLogQueryWhere[] => {
  const components = wheres.split(',')
  const result: AuditLogQueryWhere[] = []
  for (const component of components) {
    if (component === '') {
      continue
    }
    const [key, value] = component.split(':')
    switch (key) {
      case 'Action':
        result.push({
          inAction: value.split('|') as AuditLogQueryWhere['inAction'],
        })
        break
      case 'ActorId':
        result.push({
          inActorId: value.split('|') as AuditLogQueryWhere['inActorId'],
        })
        break
      case 'ActorType':
        result.push({
          inActorType: value.split('|') as AuditLogQueryWhere['inActorType'],
        })
        break
      case 'ActorOwnerId':
        result.push({
          inActorOwnerId: value.split('|') as AuditLogQueryWhere['inActorOwnerId'],
        })
        break
      case 'Id':
        result.push({
          inId: value.split('|') as AuditLogQueryWhere['inId'],
        })
        break
      case 'TargetType':
        result.push({
          inTargetType: value.split('|') as AuditLogQueryWhere['inTargetType'],
        })
        break
      case 'TargetId':
        result.push({
          inTargetId: value.split('|') as AuditLogQueryWhere['inTargetId'],
        })
        break
      case 'TargetOwnerId':
        result.push({
          inTargetOwnerId: value.split('|') as AuditLogQueryWhere['inTargetOwnerId'],
        })
        break
      case 'MinCreatedAt':
        result.push({
          minCreatedAt: value.replaceAll('_', ':'),
        })
        break
      case 'MaxCreatedAt':
        result.push({
          maxCreatedAt: value.replaceAll('_', ':'),
        })
        break
      default:
        console.warn(new Error(`Unknown where: ${JSON.stringify(key)}`))
    }
  }
  return result
}

const encodeAuditLogQueryLimit = (limit: number): string => {
  if (limit === limitDefault) {
    return ''
  }
  return `${limit}`
}

const decodeAuditLogQueryLimit = (limit: string): number => {
  if (limit === '') {
    return limitDefault
  }
  return parseInt(limit)
}

const encodeAuditLogQueryCursor = (cursor: string): string => {
  return encodeURIComponent(cursor)
}

const decodeAuditLogQueryCursor = (cursor: string): string => {
  return decodeURIComponent(cursor)
}

const sortsQP = 's'
const wheresQP = 'w'
const limitQP = 'l'
const limitDefault = 100
const cursorQP = 'c'
const pageURLBase = '/audit-logs'

export const urlReactiveAuditLogQuery = (fromQueryReactiveWithDefault: (key: string, defaultValue: string) => WritableComputedRef<string>): WritableComputedRef<AuditLogQueryReq> => {
  const qSorts = fromQueryReactiveWithDefault(sortsQP, '')
  const qWheres = fromQueryReactiveWithDefault(wheresQP, '')
  const qLimit = fromQueryReactiveWithDefault(limitQP, '')
  const qCursor = fromQueryReactiveWithDefault(cursorQP, '')

  return computed({
    get: (): AuditLogQueryReq => {
      return {
        sorts: decodeAuditLogQuerySorts(qSorts.value),
        wheres: decodeAudtLogQueryWheres(qWheres.value),
        limit: decodeAuditLogQueryLimit(qLimit.value),
        cursor: decodeAuditLogQueryCursor(qCursor.value),
      }
    },
    set: (value: AuditLogQueryReq) => {
      qSorts.value = encodeAuditLogQuerySorts(value.sorts ?? [])
      qWheres.value = encodeAuditLogQueryWheres(value.wheres)
      qLimit.value = encodeAuditLogQueryLimit(value.limit ?? limitDefault)
      qCursor.value = encodeAuditLogQueryCursor(value.cursor ?? '')
    },
  })
}

export const createURLAuditLogQuery = (localePath: LocalePathFunction, req: AuditLogQueryReq): string => {
  const qSorts = encodeAuditLogQuerySorts(req.sorts ?? [])
  const qWheres = encodeAuditLogQueryWheres(req.wheres)
  const qLimit = encodeAuditLogQueryLimit(req.limit ?? limitDefault)
  const qCursor = encodeAuditLogQueryCursor(req.cursor ?? '')
  const q = new URLSearchParams()
  if (qSorts) {
    q.set(sortsQP, qSorts)
  }
  if (qWheres) {
    q.set(wheresQP, qWheres)
  }
  if (qLimit) {
    q.set(limitQP, qLimit)
  }
  if (qCursor) {
    q.set(cursorQP, qCursor)
  }
  return localePath(pageURLBase + '?' + q.toString())
}
