/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type { AnalysisType } from './AnalysisType';

export type RunAnalysisReq = {
    /**
     * the type of analysis that should be run
     */
    analysisType: AnalysisType;
    /**
     * The pacta model that should be used to generate this analysis
     */
    pactaVersionId?: string;
    /**
     * the human meaningful name of the analysis, editable by the user
     */
    name: string;
    /**
     * Additional information about the analysis, editable by the user
     */
    description: string;
    /**
     * If populated, this analysis should be run on this portfolio
     */
    portfolioId?: string;
    /**
     * If populated, this analysis should be run on this portfolio group
     */
    portfolioGroupId?: string;
    /**
     * If populated, this analysis should be run on this initiative
     */
    initiativeId?: string;
};

