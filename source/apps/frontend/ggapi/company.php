<?php

$tokenTest = "1//0dDOKR7lbrF8hCgYIARAAGA0SNwF-L9IrZPwS2zO2ZXhsPa_fi_1oBMG_DCIqIaW4fKUIyxVUwqK-XybH_2EltGaZeRLnVhPVXQs"; // tungdtdev@gmail.com
//$tokenTest = "1//0dvT-i5eQx3eTCgYIARAAGA0SNgF-L9IrGMzF1iuU1jnr6X0XNPBHc9ahaBcQTecId3T58_EM9BMEMFLR-2l2B9M8RgULVh5tKw"; // tungdt@83.com.vn
include("config.php");

use Google\AdsApi\AdManager\Util\v202302\StatementBuilder;
use Google\AdsApi\AdManager\v202302\ServiceFactory;
use Google\AdsApi\AdManager\v202302\HasDestinationUrlCreative;
use Google\AdsApi\AdManager\v202302\BaseVideoCreative;
use Google\AdsApi\AdManager\v202302\VideoRedirectCreative;

$serviceFactory = new ServiceFactory();

$companyService = $serviceFactory->createCompanyService($session);

// Create a statement to select companies.
$pageSize = StatementBuilder::SUGGESTED_PAGE_LIMIT;
$statementBuilder = (new StatementBuilder())->orderBy('id ASC')
    ->limit($pageSize);

// Retrieve a small amount of companies at a time, paging
// through until all companies have been retrieved.
$totalResultSetSize = 0;
do {
    $page = $companyService->getCompaniesByStatement(
        $statementBuilder->toStatement()
    );

    // Print out some information for each company.
    if ($page->getResults() !== null) {
        $totalResultSetSize = $page->getTotalResultSetSize();
        $i = $page->getStartIndex();
        foreach ($page->getResults() as $company) {
            printf(
                "%d) Company with ID %d, name '%s', and type '%s' was"
                    . " found.%s",
                $i++,
                $company->getId(),
                $company->getName(),
                $company->getType(),
                PHP_EOL
            );
        }
    }

    $statementBuilder->increaseOffsetBy($pageSize);
} while ($statementBuilder->getOffset() < $totalResultSetSize);

printf("Number of results found: %d%s", $totalResultSetSize, PHP_EOL);
