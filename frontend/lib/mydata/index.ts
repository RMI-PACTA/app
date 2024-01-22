import { type LocalePathFunction } from 'vue-i18n-routing'

const pageURLBase = '/my-data'

export const QueryParamTab = 'tab'
export const QueryParamSelectedPortfolioIds = 'sp'
export const QueryParamExpandedPortfolioIds = 'ep'
export const QueryParamSelectedPortfolioGroupIds = 'sg'
export const QueryParamExpandedPortfolioGroupIds = 'eg'
export const QueryParamSelectedAnalysisIds = 'sa'
export const QueryParamExpandedAnalysisIds = 'ea'

export enum Tab {
  Portfolio = 'p',
  PortfolioGroup = 'g',
  IncompleteUpload = 'i',
  Analysis = 'a',
}

export const linkToPortfolio = (localePath: LocalePathFunction, id: string): string => {
  return linkToPortfolios(localePath, [id])
}

export const linkToPortfolios = (localePath: LocalePathFunction, ids: string[]): string => {
  const q = new URLSearchParams()
  q.set(QueryParamExpandedPortfolioIds, ids.join(','))
  q.set(QueryParamTab, Tab.Portfolio)
  return localePath(pageURLBase + '?' + q.toString())
}

export const linkToPortfolioGroup = (localePath: LocalePathFunction, id: string): string => {
  const q = new URLSearchParams()
  q.set(QueryParamExpandedPortfolioGroupIds, id)
  q.set(QueryParamTab, Tab.PortfolioGroup)
  return localePath(pageURLBase + '?' + q.toString())
}

export const linkToAnalysis = (localePath: LocalePathFunction, id: string): string => {
  const q = new URLSearchParams()
  q.set(QueryParamSelectedAnalysisIds, id)
  q.set(QueryParamTab, Tab.Analysis)
  return localePath(pageURLBase + '?' + q.toString())
}

export const linkToIncompleteUpload = (localePath: LocalePathFunction): string => {
  const q = new URLSearchParams()
  q.set(QueryParamTab, Tab.IncompleteUpload)
  return localePath(pageURLBase + '?' + q.toString())
}

export const linkToPortfolioList = (localePath: LocalePathFunction): string => {
  const q = new URLSearchParams()
  q.set(QueryParamTab, Tab.Portfolio)
  return localePath(pageURLBase + '?' + q.toString())
}
