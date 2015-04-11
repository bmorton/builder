(function(){
  var builderControllers = angular.module('builderControllers', []);

  builderControllers.controller('BuildListCtrl', [ '$scope', '$http', '$interval',
    function($scope, $http, $interval){
      $scope.builds = [];

      $scope.loadData = function(){
        $http.get('http://localhost:3000/builds').success(function(data){
          $scope.builds = data;
        });
      };

      $scope.loadData();

      $interval(function() {
        $scope.loadData();
      }, 3000);
    }]);

  builderControllers.controller('BuildDetailCtrl', ['$scope', '$routeParams',
    function($scope, $routeParams) {
      $scope.buildId = $routeParams.buildId;
      $scope.buildData = '';

      $scope.addData = function (msg) {
        $scope.$apply(function () { $scope.buildData = $scope.buildData + msg.data + "\n"; });
      };

      $scope.listen = function () {
        $scope.buildFeed = new EventSource('/builds/' + $scope.buildId);
        $scope.buildFeed.addEventListener('message', $scope.addData, false);
      };

      $scope.listen();
    }]);

})();
