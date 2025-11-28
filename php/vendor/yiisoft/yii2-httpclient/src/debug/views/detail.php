<?php
/* @var $panel yii\httpclient\debug\HttpClientPanel */
/* @var $searchModel yii\httpclient\debug\SearchModel */
/* @var $dataProvider yii\data\ArrayDataProvider */

use yii\helpers\Html;
use yii\grid\GridView;

?>
<h1><?= $panel->getName(); ?> Requests</h1>

<?php

echo GridView::widget([
    'dataProvider' => $dataProvider,
    'id' => 'db-panel-detailed-grid',
    'options' => ['class' => 'detail-grid-view table-responsive'],
    'filterModel' => $searchModel,
    'filterUrl' => $panel->getUrl(),
    'columns' => [
        ['class' => 'yii\grid\SerialColumn'],
        [
            'attribute' => 'seq',
            'label' => 'Time',
            'value' => function ($data) {
                $timeInSeconds = $data['timestamp'] / 1000;
                $millisecondsDiff = (int) (($timeInSeconds - (int) $timeInSeconds) * 1000);

                return date('H:i:s.', (int) $timeInSeconds) . sprintf('%03d', $millisecondsDiff);
            },
            'headerOptions' => [
                'class' => 'sort-numerical'
            ]
        ],
        [
            'attribute' => 'duration',
            'value' => function ($data) {
                return sprintf('%.1f ms', $data['duration']);
            },
            'options' => [
                'width' => '10%',
            ],
            'headerOptions' => [
                'class' => 'sort-numerical'
            ]
        ],
        [
            'attribute' => 'type',
            'value' => function ($data) {
                    return Html::encode($data['type']);
                },
            'filter' => $panel->getTypes(),
        ],
        [
            'attribute' => 'method',
            'value' => function ($data) {
                return Html::encode(mb_strtoupper($data['method'], 'utf8'));
            },
            'filter' => $panel->getMethods(),
        ],
        [
            'attribute' => 'request',
            'value' => function ($data) {
                $query = Html::encode($data['request']);

                if (!empty($data['trace'])) {
                    $query .= Html::ul($data['trace'], [
                        'class' => 'trace',
                        'item' => function ($trace) {
                            return "<li>{$trace['file']} ({$trace['line']})</li>";
                        },
                    ]);
                }

                if ($data['type'] !== 'batch') {
                    $query .= Html::tag(
                        'div',
                        implode('<br>', [
                            Html::a('&gt;&gt; Execute', ['request-execute', 'seq' => $data['seq'], 'tag' => Yii::$app->controller->summary['tag']], ['target' => '_blank']),
                            Html::a('&gt;&gt; Pass Through', ['request-execute', 'seq' => $data['seq'], 'tag' => Yii::$app->controller->summary['tag'], 'passthru' => true], ['target' => '_blank']),
                        ]),
                        ['class' => 'db-explain']
                    );
                }

                return $query;
            },
            'format' => 'raw',
            'options' => [
                'width' => '60%',
            ],
        ]
    ],
]);
?>
