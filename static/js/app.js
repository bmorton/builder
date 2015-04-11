(function(){
  var app = angular.module('builderApp', [
    'ngRoute',
    'builderControllers',
    'luegg.directives'
  ]);

  app.config(['$routeProvider',
    function($routeProvider) {
      $routeProvider.
        when('/builds', {
          templateUrl: 'partials/build-list.html',
          controller: 'BuildListCtrl'
        }).
        when('/builds/:buildId', {
          templateUrl: 'partials/build-detail.html',
          controller: 'BuildDetailCtrl'
        }).
        otherwise({
          redirectTo: '/builds'
        });
    }]);

})();
