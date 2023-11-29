<?php
ini_set("display_errors", 1);
header("Content-type: text/html; charset=utf-8");

if (strpos(dirname(__FILE__), 'home/assyrian/go/selfserve') > - 1) {
    //server
    define("PATH", "/home/assyrian/go/selfserve/source/www/themes/muze/assets");
    $servername = "192.168.9.10";
    $username = "apacadmin";
    $password = "iK29&6%!9XKjs@";
    $databasename = "apac_ss";
    $port = "9682";
    $encoding = "utf8mb4";
} else {
    // local
    define("PATH", "D:/ProjectGo/PubPower/source/www/themes/muze/assets");
    $servername = "localhost";
    $username = "apac_ss";
    $password = "iK29&6%!9XKjs@";
    $databasename = "apac_ss";
    $port = "3306";
    $encoding = "utf8mb4";
}

require_once PATH . '/vendor/vendor/autoload.php';

// Create connection
$conn = mysqli_connect($servername, $username, $password, $databasename, $port, $encoding);
mysqli_set_charset($conn, 'UTF8');
// Check connection
if (!$conn) {
    echo json_encode(["status" => FALSE, "message" => "Build invoice PDF: Connect mysql error"]);
    return;
}
$sql = "SELECT * FROM payment_invoice WHERE (pdf = '' and status = 2) OR (status = 1 and statement = '')";
$result = $conn->query($sql);
if (!empty($result) && $result->num_rows > 0) {
    while ($row = $result->fetch_assoc()) {
        $createPDF = createPDF($row, $conn);
    }
}

echo json_encode(["status" => TRUE, "message" => "Scan PDF success!"]);
return;

function createPDF($invoice, $conn) {
    if (!$invoice) {
        return TRUE;
    }

    $user = $billing = $requests = [];
    $sql = "SELECT * FROM user WHERE id = " . $invoice['user_id'];
    $result = $conn->query($sql);
    if (!empty($result) && $result->num_rows > 0) {
        while ($row = $result->fetch_assoc()) {
            $user = $row;
        }
    }
    if (!$user) {
        $sql = "INSERT INTO notification (user_id, status, message, action, link, created_at) VALUES (7, 1, 'Build invoice PDF: " . $invoice['user_id'] . " User does not exist!', '','' , '" . date('Y-m-d H-i-s') . "')";
        $result = $conn->query($sql);
        return TRUE;
        // return json_encode(["status" => FALSE, "message" => $invoice['user_id'] . " User does not exist!"]);
    }

    $sql = "SELECT * FROM user_billing WHERE user_id = " . $user['id'];
    $result = $conn->query($sql);
    if (!empty($result) && $result->num_rows > 0) {
        while ($row = $result->fetch_assoc()) {
            $billing = $row;
        }
    }

    $sql = "SELECT * FROM payment_request WHERE id in (" . $invoice['request_id'] . ")";
    $result = $conn->query($sql);
    if (!empty($result) && $result->num_rows > 0) {
        while ($row = $result->fetch_assoc()) {
            $requests[] = $row;
        }
    }
    if (!$requests) {
        $sql = "INSERT INTO notification (user_id, status, message, action, link, created_at) VALUES (7, 1, 'Build invoice PDF: " . $invoice['id'] . " not requests!', '','' , '" . date('Y-m-d H-i-s') . "')";
        $result = $conn->query($sql);
        // return json_encode(["status" => FALSE, "message" => $invoice['id'] . " not requests!"]);
        return TRUE;
    }

    if ($invoice['status'] == 2) {
        $result_pdf = renderInvoicePDF($invoice, $billing, $user, $requests, $conn);
    } else {
        $result_pdf = renderStatementPDF($invoice, $user, $requests, $conn);
    }
    return $result_pdf;
}

function renderStatementPDF($invoice, $user, $requests, $conn) {
    $total_amount_request = 0;
    foreach ($requests as $request) {
        $total_amount_request += $request['amount'];
    }
    $config = [
        'mode'             => 'utf-8',
        'format'           => 'A4',
        'default_font'     => 'FreeSerif',
        'orientation'      => 'P',
        'debug'            => TRUE,
        'autoScriptToLang' => TRUE,
        "autoLangToFont"   => TRUE,
    ];
    $mpdf = new \Mpdf\Mpdf($config);

    ob_start();
    include("statement.php");
    $HTML = ob_get_clean();
    // $HTML =  "<div>reprot Thiếu từ</div>";
    utf8_encode($HTML);   //chỉ chuyển đổi ISO-8859-1 thành UTF-8
    // $HTML = iconv('UTF-8', 'UTF-8', $HTML);
    // $HTML = mb_convert_encoding($HTML, 'UTF-8', 'UTF-8');
    // $HTML = iconv('windows-1252', 'UTF-8', $HTML);
    // $HTML = iconv('utf-8', 'us-ascii//TRANSLIT', $HTML);
    // var_dump( mb_detect_encoding($HTML, mb_detect_order(), TRUE) );die;
    $HTML = iconv(mb_detect_encoding($HTML, mb_detect_order(), TRUE), "UTF-8", $HTML);
    $mpdf->allow_charset_conversion = TRUE;  // Set by default to TRUE
    $mpdf->charset_in = 'UTF-8';

    $mpdf->WriteHTML($HTML);
    $result = $mpdf->Output(PATH . "/invoice/PP" . $invoice["id"] . '-statement-' . $invoice['start_date'] . ".pdf");

    $sql = "UPDATE payment_invoice SET statement = '/assets/invoice/PP" . $invoice["id"] . "-statement-" . $invoice['start_date'] . ".pdf' WHERE id = " . $invoice['id'];
    $result = $conn->query($sql);
    echo json_encode(["status" => TRUE, "path" => "/assets/invoice/PP" . $invoice["id"] . '-statement-' . $invoice['start_date'] . ".pdf"]);
    return TRUE;
}

function renderInvoicePDF($invoice, $billing, $user, $requests, $conn) {
    $total_amount_request = 0;
    foreach ($requests as $request) {
        $total_amount_request += $request['amount'];
    }
    $config = [
        'mode'             => 'utf-8',
        'format'           => 'A4',
        'default_font'     => 'FreeSerif',
        'orientation'      => 'P',
        'debug'            => TRUE,
        'autoScriptToLang' => TRUE,
        "autoLangToFont"   => TRUE,
    ];
    $mpdf = new \Mpdf\Mpdf($config);

    ob_start();
    include("invoice.php");
    $HTML = ob_get_clean();
    // $HTML =  "<div>reprot Thiếu từ</div>";
    utf8_encode($HTML);   //chỉ chuyển đổi ISO-8859-1 thành UTF-8
    // $HTML = iconv('UTF-8', 'UTF-8', $HTML);
    // $HTML = mb_convert_encoding($HTML, 'UTF-8', 'UTF-8');
    // $HTML = iconv('windows-1252', 'UTF-8', $HTML);
    // $HTML = iconv('utf-8', 'us-ascii//TRANSLIT', $HTML);
    // var_dump( mb_detect_encoding($HTML, mb_detect_order(), TRUE) );die;
    $HTML = iconv(mb_detect_encoding($HTML, mb_detect_order(), TRUE), "UTF-8", $HTML);
    $mpdf->allow_charset_conversion = TRUE;  // Set by default to TRUE
    $mpdf->charset_in = 'UTF-8';
    // $mpdf->charset_in = 'UTF-8';

    $mpdf->WriteHTML($HTML);
    $mpdf->Output(PATH . "/invoice/PP" . $invoice["id"] . '-' . $invoice['start_date'] . ".pdf");

    $sql = "UPDATE payment_invoice SET pdf = '/assets/invoice/PP" . $invoice["id"] . "-" . $invoice['start_date'] . ".pdf' WHERE id = " . $invoice['id'];
    $result = $conn->query($sql);
    echo json_encode(["status" => TRUE, "path" => "/assets/invoice/PP" . $invoice["id"] . '-' . $invoice['start_date'] . ".pdf"]);
    return TRUE;
}