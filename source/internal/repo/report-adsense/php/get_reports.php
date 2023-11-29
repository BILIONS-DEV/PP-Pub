<?php

include("config.php");

$val = getopt(null, ["account:","startDate:","endDate:"]);

if (empty($val["account"])) {
    echo json_encode(["status" => false, "message" => "account is required"]);
    return;
}
$account = '';
if ($val["account"]) {
    $account = $val["account"];
}

if (empty($val["startDate"])) {
    echo json_encode(["status" => false, "message" => "startDate is required"]);
    return;
}
$startDate = '';
if ($val["startDate"]) {
    $startDate = $val["startDate"];
}

if (empty($val["endDate"])) {
    echo json_encode(["status" => false, "message" => "endDate is required"]);
    return;
}
$endDate = '';
if ($val["endDate"]) {
    $endDate = $val["endDate"];
}
list($yearStartDate, $monthStartDate, $dayStartDate) = explode("-", $startDate);
list($yearEndDate, $monthEndDate, $dayEndDate) = explode("-", $endDate);
$optParams = array(
    'startDate.year' => $yearStartDate,
    'startDate.month' => $monthStartDate,
    'startDate.day' => $dayStartDate,
    'endDate.year' => $yearEndDate,
    'endDate.month' => $monthEndDate,
    'endDate.day' => $dayEndDate,
    'currencyCode' => 'USD',
//    'reportingTimeZone' => 'GOOGLE_TIME_ZONE',
    'metrics' => array(
        'PAGE_VIEWS', 'IMPRESSIONS', 'CLICKS',
        'AD_REQUESTS_CTR', 'COST_PER_CLICK', 'AD_REQUESTS_RPM',
        'ESTIMATED_EARNINGS'),
    'dimensions' => array(
        'date', 'CUSTOM_CHANNEL_ID', 'CUSTOM_CHANNEL_NAME', 'COUNTRY_CODE'),
//      'orderBy' => '+DIMENSION_UNSPECIFIED',
    'filters' => array(//        'DATE=='.$date
    )
);

// Run report.
$report = $service->accounts_reports->generate($account, $optParams);

$results = [];
if (isset($report) && isset($report['rows'])) {
    // Display results.
    foreach ($report['rows'] as $row) {
        $result = [];
        foreach ($row['cells'] as $indexCell => $column) {
            $result[strtolower($report['headers'][$indexCell]["name"])] = $column["value"];
        }
        array_push($results, $result);
    }
}
echo json_encode(["status" => true, "type" => "GET_REPORTS", "message" => "Success!", "data" => $results]);