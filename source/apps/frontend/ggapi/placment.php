<?php

$tokenTest = "1//0dFAe9Ha9dDsSCgYIARAAGA0SNwF-L9IrUcewgNvSsrj5EZTAqFRx7VZ1AWeZP4zQXOm6QValz3jArpoZ9boH178Kwd2Cdabu1uo";
include("config.php");

use Google\AdsApi\AdManager\Util\v202302\StatementBuilder;
use Google\AdsApi\AdManager\v202302\ServiceFactory;

$serviceFactory = new ServiceFactory();
$placementService = $serviceFactory->createPlacementService($session);

$pageSize = StatementBuilder::SUGGESTED_PAGE_LIMIT;
$statementBuilder = (new StatementBuilder())
    ->orderBy('id ASC')
    ->limit($pageSize);

// Retrieve a small amount of placements at a time, paging
// through until all placements have been retrieved.
$totalResultSetSize = 0;
do {
    $page = $placementService->getPlacementsByStatement(
        $statementBuilder->toStatement()
    );

    // Print out some information for each placement.
    if ($page->getResults() !== null) {
        $totalResultSetSize = $page->getTotalResultSetSize();
        $i = $page->getStartIndex();
        foreach ($page->getResults() as $placement) {
            printf(
                "%d) Placement with ID %d and name '%s' was found.%s",
                $i++,
                $placement->getId(),
                $placement->getName(),
                PHP_EOL
            );
        }
    }

    $statementBuilder->increaseOffsetBy($pageSize);
} while ($statementBuilder->getOffset() < $totalResultSetSize);

printf("Number of results found: %d%s", $totalResultSetSize, PHP_EOL);
