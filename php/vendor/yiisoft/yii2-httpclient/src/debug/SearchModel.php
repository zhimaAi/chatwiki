<?php
/**
 * @link https://www.yiiframework.com/
 * @copyright Copyright (c) 2008 Yii Software LLC
 * @license https://www.yiiframework.com/license/
 */

namespace yii\httpclient\debug;

use yii\data\ArrayDataProvider;
use yii\debug\components\search\Filter;
use yii\debug\models\search\Base;

/**
 * Log search model.
 *
 * @author Paul Klimov <klimov.paul@gmail.com>
 * @since 2.0
 */
class SearchModel extends Base
{
    /**
     * @var string type of the input search value
     */
    public $type;
    /**
     * @var string method of the input search value
     */
    public $method;
    /**
     * @var int request attribute input search value
     */
    public $request;


    /**
     * {@inheritdoc}
     */
    public function rules()
    {
        return [
            [['type', 'method', 'request'], 'safe'],
        ];
    }

    /**
     * {@inheritdoc}
     */
    public function attributeLabels()
    {
        return [
            'type' => 'Type',
            'method' => 'Method',
            'request' => 'Request',
        ];
    }

    /**
     * Returns data provider with filled models. Filter applied if needed.
     *
     * @param array $params an array of parameter values indexed by parameter names
     * @param array $models data to return provider for
     * @return \yii\data\ArrayDataProvider
     */
    public function search($params, $models)
    {
        $dataProvider = new ArrayDataProvider([
            'allModels' => $models,
            'pagination' => false,
            'sort' => [
                'attributes' => ['duration', 'seq', 'type', 'method', 'request'],
                'defaultOrder' => [
                    'duration' => SORT_DESC,
                ],
            ],
        ]);

        if (!($this->load($params) && $this->validate())) {
            return $dataProvider;
        }

        $filter = new Filter();
        $this->addCondition($filter, 'type', true);
        $this->addCondition($filter, 'method', true);
        $this->addCondition($filter, 'request', true);
        $dataProvider->allModels = $filter->filter($models);

        return $dataProvider;
    }
}