<?php

use Google\AdsApi\AdManager\AdManagerSession;
use Google\AdsApi\AdManager\v202302\ServiceFactory;
use Google\AdsApi\AdManager\Util\v202302\StatementBuilder;

include("config_api.php");

// Check name truyền từ command
if (empty($val["name"])) {
    echo json_encode(["status" => false, "message" => "name is required"]);
    return;
} else {
    $name = $val["name"];
}
getAdUnit($name, $dev, $session);

function getAdUnit($name, $dev, AdManagerSession $session) {
    $serviceFactory = new ServiceFactory();
    $inventoryService = $serviceFactory->createInventoryService($session);

    // Create a statement to select ad units.
    $pageSize = StatementBuilder::SUGGESTED_PAGE_LIMIT;
    $statementBuilder = (new StatementBuilder())->where('name = :name')
        ->orderBy('id ASC')
        ->limit(1)
        ->withBindVariableValue('name', $name);

    // Retrieve a small amount of ad units at a time, paging
    // through until all ad units have been retrieved.
    $totalResultSetSize = 0;
    $adUnitId = 0;
    do {
        try {
            $page = $inventoryService->getAdUnitsByStatement(
                $statementBuilder->toStatement()
            );
        } catch (\Google\AdsApi\AdManager\v202302\ApiException $e) {
            echo json_encode(["status" => false, "type" => "adUnit", "message" => "Failed to get adUnit!", "log" => $inventoryService->getLastSoapFaultMessage()], JSON_NUMERIC_CHECK);
            return;
        }

        // Print out some information for each ad unit.
        if ($page->getResults() !== null) {
            $totalResultSetSize = $page->getTotalResultSetSize();
            $i = $page->getStartIndex();
            foreach ($page->getResults() as $adUnit) {
                $adUnitId = $adUnit->getId();
                if ($dev) {
                    printf(
                        "%d) Ad unit with ID '%s' and name '%s' was found.%s",
                        $i++,
                        $adUnit->getId(),
                        $adUnit->getName(),
                        PHP_EOL
                    );
                }
            }
        }
        $statementBuilder->increaseOffsetBy($pageSize);
    } while ($statementBuilder->getOffset() < $totalResultSetSize);
    if ($adUnitId != 0) {
        echo json_encode(["status" => true, "type" => "adUnit", "message" => "Success!", "data" => strval($adUnitId)]);
    } else {
        echo json_encode(["status" => false, "type" => "adUnit", "message" => "Not found!"]);
    }
}
