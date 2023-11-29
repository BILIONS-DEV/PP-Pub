<?php
$showRate = FALSE;
if ($requests) {
    foreach ($requests as $request) {
        if (!empty($request['rate'])) {
            $showRate = TRUE;
        }
    }
}
?>
<html lang="en">
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
</head>
<body>
<div class="container-fluid px-0">
    <div class="row">
        <div class="col-12">
            <div class="d-flex flex-wrap">
                <span class="mb-3 mb-md-0 me-1">
                    <img src="https://pubpower.io/img/pub-power-black.png" alt="Pubpower" width="200" class="dark-logo">
                </span>
            </div>
            <div class="d-flex justify-content-between" style="display: inline">
                <div class=" pt-5 pb-4 pull-left" style="width: 50%">
                    <h5 class="font-weight-semibold opensans-font"><strong>From: APAC Digital Exchange</strong></h5>
                    <div class="text-gray-700 mb-0">Vistra Singapore 9 Raffles Place #26-01, Republic Plaza Singapore 048619</div>
                </div>
                <div class=" pt-5 pb-4 pull-right" style="width: 50%;text-align: right">
                    <h5 class="font-weight-semibold opensans-font">
                        <strong>Billed to: <?= $user['email'] ?></strong>
                    </h5>
                    <?php if (!empty($user['address'])) {
                        echo '<div class="text-gray-700 mb-0">' . $user['address'] . (!empty($user['city']) ? ', ' . $user['city'] : '') . (!empty($user['country']) ? ', ' . $user['country'] : '') . '</div>';
                    } ?>
                </div>
            </div>
            <div class="border-top border-gray-200 pt-4 mt-4">
                <div class="table-responsive">
                    <table class="table table-nowrap">
                        <thead>
                        <tr>
                            <th>Invoice ID</th>
                            <th>Paid Date</th>
                            <th>Issue Date</th>
                            <th>Status</th>
                        </tr>
                        </thead>
                        <tbody>
                        <tr>
                            <td>#PP<?= $invoice['id'] ?></td>
                            <td><?= date("m/d/y", strtotime($invoice['paid_date'])) ?></td>
                            <td><?= date("m/d/y", strtotime($invoice['end_date'])) ?> </td>
                            <td><?= $invoice['status'] == 2 ? 'Paid' : 'Pending' ?></td>
                        </tr>
                        </tbody>
                    </table>
                </div>
            </div>
            <div class="border-top border-gray-200 pt-5 mt-4 mb-4">
                <div class="table-responsive">
                    <table class="table border table-bordered card-table table-nowrap">
                        <thead>
                        <tr>
                            <th>Period</th>
                            <th>Type</th>
                            <th>Note</th>
                             <?php if (!empty($showRate)): ?>
                                <th>Revenue</th>
                                <th>Rate</th>
                            <?php endif; ?>
                            <th>Amount</th>
                        </tr>
                        </thead>
                        <tbody>
                        <?php if ($requests) {
                            foreach ($requests as $request) { ?>
                                <tr>
                                    <td><?= date("m/d/Y", strtotime($request['start_date'])) . "  " . date("m/d/Y", strtotime($request['end_date'])) ?></td>
                                    <td><?= $request['type'] == 1 ? 'Commission' : 'Prepaid' ?></td>
                                    <td><?= $request['note'] ?></td>
                                    <?php if (!empty($showRate)): ?>
                                       <td>$<?= number_format($request['revenue'], 2) ?></th>
                                       <td><?= $request['rate'] ?>%</th>
                                    <?php endif; ?>
                                    <td>$<?= number_format($request['amount'], 2) ?></td>
                                </tr>
                            <?php } ?>
                        <?php } ?>
                        </tbody>
                    </table>
                </div>
                <div class="pull-right" style="width: 30%;">
                    <div class="">
                        <table class="table table-total table-nowrap">
                            <tbody>
                               <?php if (empty($showRate)) { ?>
                                    <tr>
                                        <td>Subtotal:</td>
                                        <td>$<?= number_format($total_amount_request, 2) ?></td>
                                    </tr>
                                    <tr>
                                        <td>Rate:</td>
                                       <td><?= !empty($invoice['rate']) ? $invoice['rate'] : 100 ?>%</td>
                                    </tr>
                               <?php } ?>
                            <tr>
                                <td><strong>Total:</strong></td>
                                <td><strong>$<?= number_format($invoice['amount'], 2) ?></strong></td>
                            </tr>
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
            <div class="border-top border-gray-200 pt-4 mt-4">
                <div class="d-flex justify-content-between" style="display: inline">
                    <div class="pull-left" style="width: 100%;">
                        <div class="mt-2 pt-1">
                            <?php if (!empty($billing)) { ?>
                                <?php if ($billing['method'] == 'bank') { ?>
                                    <div class="">
                                        <table class="table table-total table-nowrap">
                                            <tbody>
                                            <tr>
                                                <td>
                                                    <i class="fa fa-university text-success mr-2" style="font-size: 20px"></i>
                                                    <strong>Bank Transfer</strong>
                                                </td>
                                                <td></td>
                                            </tr>
                                            <tr>
                                                <td>Beneficiary Name:</td>
                                                <td class="text-end"><?= $billing['beneficiary_name'] ?></td>
                                            </tr>
                                            <tr>
                                                <td>Bank Name:</td>
                                                <td class="text-end"><?= $billing['bank_name'] ?></td>
                                            </tr>
                                            <tr>
                                                <td>Bank Address:</td>
                                                <td class="text-end"><?= $billing['bank_address'] ?></td>
                                            </tr>
                                            <tr>
                                                <td>Bank Account Number:</td>
                                                <td class="text-end"><?= $billing['bank_account_number'] ?></td>
                                            </tr>
                                            <tr>
                                                <td>Bank Routing Number:</td>
                                                <td class="text-end"><?= $billing['bank_routing_number'] ?></td>
                                            </tr>
                                            <tr>
                                                <td>Bank Iban Number:</td>
                                                <td class="text-end"><?= $billing['bank_iban_number'] ?></td>
                                            </tr>
                                            <tr>
                                                <td>Swift Code:</td>
                                                <td class="text-end"><?= $billing['swift_code'] ?></td>
                                            </tr>
                                            </tbody>
                                        </table>
                                    </div>
                                <?php } else if ($billing['method'] == 'payoneer') { ?>
                                    <div class="row mt-md-2 py-2 py-md-2 pe-md-3">
                                        <div class="col-5">
                                            <img src="https://apps.valueimpression.com/assets/img/payoneer.png" alt="" width="50px">
                                        </div>
                                        <div class="col-5 col-sm-3 col-xxl-2 " style="line-height: 25px;">
                                            <span class="text-black-600 text-nowrap"><?= $billing['payoneer_email'] ?></span>
                                        </div>
                                    </div>
                                <?php } else if ($billing['method'] == 'paypal') { ?>
                                    <div class="row mt-md-2 py-2 py-md-2 pe-md-3">
                                        <div class="col-5">
                                            <img src="https://apps.valueimpression.com/assets/img/paypal.png" alt="" width="50px">
                                        </div>
                                        <div class="col-5 col-sm-3 col-xxl-2 " style="line-height: 30px;">
                                            <span class="text-black-600 text-nowrap"><?= $billing['paypal_email'] ?></span>
                                        </div>
                                    </div>
                                <?php } else if ($billing['method'] == 'currency') { ?>
                                    <div class="row mt-md-2 py-2 py-md-2 pe-md-3">
                                        <div class="col-5">
                                            <span class="font-weight-semibold text-black-600"><?= $billing['cryptocurrency'] ?>:</span>
                                        </div>
                                        <div class="col-5 col-sm-3 col-xxl-2 ">
                                            <span class="text-black-600 text-nowrap"><?= $billing['wallet_id'] ?></span>
                                        </div>
                                    </div>
                                <?php } else { ?>
                                    <div class="row mt-md-2 py-2 py-md-2 pe-md-3">Update Billing</div>
                                <?php } ?>
                            <?php } else { ?>
                                <div class="row mt-md-2 py-2 py-md-2 pe-md-3">Update Billing</div>
                            <?php } ?>
                        </div>
                    </div>

                </div>
            </div>
        </div>
    </div>
</div>
</body>
<style>
    .alert, .brand, .btn-simple, .h1, .h2, .h3, .h4, .h5, .h6, .navbar, .td-name, a, body, button.close, h1, h2, h3, h4, h5, h6, p, td {
        -moz-osx-font-smoothing: grayscale;
        -webkit-font-smoothing: antialiased;
        font-family: 'Roboto', serif !important;
    }
    .container, .container-fluid, .container-lg, .container-md, .container-sm, .container-xl, .container-xxl {
        width: 100%;
        font-size: 12px;
        margin-right: auto;
        margin-left: auto;
    }
    .row {
        --bs-gutter-x: 1.5rem;
        --bs-gutter-y: 0;
        display: -ms-flexbox;
        display: flex;
        -ms-flex-wrap: wrap;
        flex-wrap: wrap;
        margin-top: calc(var(--bs-gutter-y) * -1);
        margin-right: calc(var(--bs-gutter-x) * -.5);
        margin-left: calc(var(--bs-gutter-x) * -.5);
    }
    .row {
        display: -webkit-box; /* wkhtmltopdf uses this one */
        display: flex;
        -webkit-box-pack: center; /* wkhtmltopdf uses this one */
        justify-content: center;
    }
    .row > div {
        -webkit-box-flex: 1;
        -webkit-flex: 1;
        flex: 1;
    }
    .row > div:last-child {
        margin-right: 0;
    }
    .col-12 {
        -ms-flex: 0 0 auto;
        flex: 0 0 auto;
        width: 100%;
    }

    p {
        line-height: 1.75rem;
    }
    .h1, .h2, .h3, .h4, .h5, .h6, h1, h2, h3, h4, h5, h6 {
        margin-top: 0;
        margin-bottom: 0.5rem;
        font-weight: 500;
        line-height: 1.2;
    }
    .col, .col-1, .col-10, .col-11, .col-12, .col-2, .col-3, .col-4, .col-5, .col-6, .col-7, .col-8, .col-9, .col-auto, .col-lg, .col-lg-1, .col-lg-10, .col-lg-11, .col-lg-12, .col-lg-2, .col-lg-3, .col-lg-4, .col-lg-5, .col-lg-6, .col-lg-7, .col-lg-8, .col-lg-9, .col-lg-auto, .col-md, .col-md-1, .col-md-10, .col-md-11, .col-md-12, .col-md-2, .col-md-3, .col-md-4, .col-md-5, .col-md-6, .col-md-7, .col-md-8, .col-md-9, .col-md-auto, .col-sm, .col-sm-1, .col-sm-10, .col-sm-11, .col-sm-12, .col-sm-2, .col-sm-3, .col-sm-4, .col-sm-5, .col-sm-6, .col-sm-7, .col-sm-8, .col-sm-9, .col-sm-auto, .col-xl, .col-xl-1, .col-xl-10, .col-xl-11, .col-xl-12, .col-xl-2, .col-xl-3, .col-xl-4, .col-xl-5, .col-xl-6, .col-xl-7, .col-xl-8, .col-xl-9, .col-xl-auto {
        position: relative;
        width: 100%;
        padding-right: 15px;
        padding-left: 15px;
    }
    .border-gray-50 {
        border-color: #e9ecef !important;
    }

    .border {
        border: 1px solid #dee2e6 !important;
    }
    .card-body {
        -ms-flex: 1 1 auto;
        flex: 1 1 auto;
        padding: 1rem 1rem;
    }
    .align-items-center {
        -ms-flex-align: center !important;
        align-items: center !important;
        text-align: center !important;
    }
    .justify-content-between {
        -ms-flex-pack: justify !important;
        justify-content: space-between !important;
    }
    .border-top {
        border-top: 1px solid #dee2e6 !important;
    }
    .col-auto {
        -ms-flex: 0 0 auto;
        flex: 0 0 auto;
        width: auto;
    }
    .font-weight-semibold {
        font-weight: 600 !important;
    }
    .avatar-border, .avatar-border-lg {
        border: 0.1875rem solid #fff;
    }
    .avatar-sm-status {
        bottom: -0.2625rem;
        right: -0.2625rem;
        width: 0.75rem;
        min-width: 0.75rem;
        height: 0.75rem;
        font-size: 12px;
    }
    .avatar-warning {
        background-color: #fd7e14;
    }
    .avatar {
        position: relative;
        display: inline-block;
        width: 3.125rem;
        height: 3.125rem;
        min-width: 3.125rem;
        border-radius: 0.375rem;
    }
    .h5, h5 {
        font-size: 13px;
        font-weight: bold;
    }

    .border-gray-200 {
        border-color: #e9ecef !important;
    }
    .d-flex {
        display: -ms-flexbox !important;
        display: flex !important;
    }
    .flex-wrap {
        -ms-flex-wrap: wrap !important;
        flex-wrap: wrap !important;
    }

    .table-responsive {
        position: relative;
        z-index: 0;
    }
    .table {
        --bs-table-bg: transparent;
        --bs-table-accent-bg: transparent;
        --bs-table-striped-color: #212529;
        --bs-table-striped-bg: rgba(0, 0, 0, 0.05);
        --bs-table-active-color: #212529;
        --bs-table-active-bg: rgba(0, 0, 0, 0.1);
        --bs-table-hover-color: #212529;
        --bs-table-hover-bg: rgba(0, 0, 0, 0.075);
        width: 100%;
        margin-bottom: 1rem;
        color: #212529;
        vertical-align: top;
        border-color: #dbdcdd;
    }
    .card-table th {
        font-size: 12px;
        text-transform: uppercase;
        font-weight: 400;
        padding: 0.8125rem 2rem;
        background-color: #dbdcdd;
        border-color: #dbdcdd !important;
    }
    .table th {
        font-size: 12px;
        font-weight: 400;
        padding: 0.5125rem 2rem;
    }
    .table td {
        font-size: 12px;
        font-weight: 400;
        padding: 0.5125rem 2rem;
    }
    .table th {
        text-align: left !important;
    }
    .card-table td {
        font-size: 12px;
        padding: 0.9375rem 2rem;
        vertical-align: middle;
        border-color: #e9ecef !important;
        color: #1e1e1e;
    }
    .pull-left {
        float: left;
    }
    .pull-right {
        float: right;
    }
     .pt-5 {
        padding-top: 3rem !important;
    }
     .mt-5 {
        margin-top: 3rem !important;
    }
     .pt-4 {
        padding-top: 1.5rem !important;
    }
    .mt-4 {
        margin-top: 1.5rem !important;
    }
    .pb-4 {
        padding-bottom: 1.5rem !important;
    }
    .mb-4 {
        margin-bottom: 1.5rem !important;
    }
    .pt-3 {
        padding-top: 1rem !important;
    }
   .p-3 {
        padding: 1rem !important;
    }
    .mt-3 {
        margin-top: 1rem !important;
    }
    .mb-2 {
        margin-bottom: 0.5rem !important;
    }
    .mt-2 {
        margin-top: 0.5rem !important;
    }
    .pb-2 {
        padding-bottom: 0.5rem !important;
    }
    .py-2{
        padding: 0.5rem 0rem !important;
    }
     .mt-2{
        margin-top: 0.5rem;
    }
     .pt-1{
        padding-top: 0.2rem !important;
    }
    .table-total td {
        font-size: 12px;
        padding: 0.3375rem;
        vertical-align: middle;
        border-color: #e9ecef!important;
        color: #1e1e1e;
    }
   .text-end {
        text-align: right !important;
    }
</style>
</html>