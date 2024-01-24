import {
  Tab,
  QueryParamTab,
  QueryParamSelectedPortfolioIds,
  QueryParamExpandedPortfolioIds,
  QueryParamExpandedPortfolioGroupIds,
  QueryParamSelectedAnalysisIds,
} from '@/lib/mydata'

const pageURLBase = '/my-data'

export const useMyDataURLs = () => {
  const localePath = useLocalePath()

  const toURL = (q: URLSearchParams): string => localePath(pageURLBase + '?' + q.toString())

  const linkToPortfolios = (ids: string[]): string => {
    const q = new URLSearchParams()
    q.set(QueryParamSelectedPortfolioIds, ids.join(','))
    q.set(QueryParamTab, Tab.Portfolio)
    return toURL(q)
  }

  const linkToPortfolio = (id: string): string => {
    const q = new URLSearchParams()
    q.set(QueryParamExpandedPortfolioIds, id)
    q.set(QueryParamTab, Tab.Portfolio)
    return toURL(q)
  }

  const linkToPortfolioGroup = (id: string): string => {
    const q = new URLSearchParams()
    q.set(QueryParamExpandedPortfolioGroupIds, id)
    q.set(QueryParamTab, Tab.PortfolioGroup)
    return toURL(q)
  }

  const linkToAnalysis = (id: string): string => {
    const q = new URLSearchParams()
    q.set(QueryParamSelectedAnalysisIds, id)
    q.set(QueryParamTab, Tab.Analysis)
    return toURL(q)
  }

  const linkToIncompleteUploadList = (): string => {
    const q = new URLSearchParams()
    q.set(QueryParamTab, Tab.IncompleteUpload)
    return toURL(q)
  }

  const linkToPortfolioList = (): string => {
    const q = new URLSearchParams()
    q.set(QueryParamTab, Tab.Portfolio)
    return toURL(q)
  }

  return {
    linkToPortfolios,
    linkToPortfolio,
    linkToPortfolioGroup,
    linkToAnalysis,
    linkToIncompleteUploadList,
    linkToPortfolioList,
  }
}
