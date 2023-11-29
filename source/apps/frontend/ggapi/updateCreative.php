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
$creativeService = $serviceFactory->createCreativeService($session);

$creativeId = "138355937637";

// Create a statement to select a single creative by ID.
$statementBuilder = (new StatementBuilder())->where('id = :id')
    ->orderBy('id ASC')
    ->limit(1)
    ->withBindVariableValue('id', $creativeId);

// Get the creative.
$page = $creativeService->getCreativesByStatement(
    $statementBuilder->toStatement()
);

$creative = $page->getResults()[0];
var_dump($creative->getName());
// Only update the destination URL if it has one.
if ($creative instanceof HasDestinationUrlCreative) {
    // Update the destination URL of the creative.
    $creative->setDestinationUrl('https://news.google.com');

    // Update the creative on the server.
    //    $creatives = $creativeService->updateCreatives([$creative]);

    foreach ($creatives as $updatedCreative) {
        printf(
            "Creative with ID %d and name '%s' was updated.%s",
            $updatedCreative->getId(),
            $updatedCreative->getName(),
            PHP_EOL
        );
    }
} else {
    printf("No creatives were updated.%s", PHP_EOL);
}
